import cv2
import threading
import time
import socket
import struct
from flask import Flask, Response, jsonify
from ultralytics import YOLO

# ================= 配置区域 =================
# 请确保这个 IP 是你 ESP32 串口监视器里打印出来的 IP
ESP32_STREAM_URL = "http://172.20.10.3/stream"
MODEL_PATH = "esp32_face_640_best.pt"

TARGET_W, TARGET_H = 160, 120
# ===========================================

app = Flask(__name__)

# 全局变量
latest_detection_result = {"status": "idle", "message": "Waiting for ESP32..."}
processed_frame = None
frame_lock = threading.Lock()

# 1. 加载 AI 模型
print(">>> 正在加载 AI 模型...")
try:
    model = YOLO(MODEL_PATH)
    print(">>> 模型加载成功！")
except Exception as e:
    print(f">>> 模型加载失败: {e}")
    exit()


# 2. AI 推理线程
def ai_processing_loop():
    global latest_detection_result, processed_frame

    # 增加重试机制，防止 ESP32 还没启动导致连接失败
    max_retries = 5
    retry_count = 0

    while retry_count < max_retries:
        print(f">>> 正在尝试连接 ESP32: {ESP32_STREAM_URL} (尝试 {retry_count + 1}/{max_retries})...")
        cap = cv2.VideoCapture(ESP32_STREAM_URL)

        if cap.isOpened():
            print(">>> ✅ 成功连接到 ESP32 视频流！开始 AI 推理...")
            retry_count = 0  # 重置计数器

            while True:
                ret, frame = cap.read()
                if not ret:
                    print(">>> ⚠️ 视频流断开，正在重连...")
                    break  # 跳出内层循环，尝试重连

                # --- AI 推理 ---
                results = model(frame, verbose=False, conf=0.6)

                # --- 绘图 ---
                annotated_frame = results[0].plot()

                # --- 更新数据 ---
                with frame_lock:
                    processed_frame = annotated_frame

                # --- 更新状态 ---
                detected_faces = results[0].boxes
                if len(detected_faces) > 0:
                    # 获取第一个检测到的人脸
                    name = results[0].names[int(detected_faces[0].cls[0])]
                    conf = float(detected_faces[0].conf[0])
                    latest_detection_result = {"status": "found", "name": name, "confidence": round(conf, 2)}
                else:
                    latest_detection_result = {"status": "empty", "message": "No face"}
        else:
            print(">>> ❌ 连接失败，请检查 ESP32 IP 地址是否正确，或 ESP32 是否已启动。")
            retry_count += 1
            time.sleep(2)  # 等待 2 秒再试

    print(">>> ⛔ 连接 ESP32 失败次数过多，AI 线程停止。")


# 3. 视频流路由 (MJPEG)
def generate_frames():
    while True:
        with frame_lock:
            if processed_frame is None:
                # 如果没有画面，发送一张黑图或提示图，防止浏览器断开
                # 这里简单返回一个空的帧，实际项目中可以放一张 "Loading..." 图片
                continue
            ret, buffer = cv2.imencode('.jpg', processed_frame)
            if not ret:
                continue

        yield (b'--frame\r\n'
               b'Content-Type: image/jpeg\r\n\r\n' + buffer.tobytes() + b'\r\n')


@app.route('/video_feed')
def video_feed():
    return Response(generate_frames(), mimetype='multipart/x-mixed-replace; boundary=frame')


# 4. API 路由 (给 ESP32 用)
@app.route('/api/status')
def get_status():
    return jsonify(latest_detection_result)


# 5. 主页路由 (优化了 HTML 写法)
@app.route('/')
def index():
    html_content = """
    <!DOCTYPE html>
    <html lang="zh-CN">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>智能门禁监控</title>
        <style>
            body { background-color: #1a1a1a; color: #fff; font-family: 'Segoe UI', sans-serif; text-align: center; margin: 0; padding: 20px; }
            h1 { margin-bottom: 10px; }
            .container { margin: 0 auto; width: 640px; max-width: 100%; position: relative; }
            img { width: 100%; border-radius: 12px; box-shadow: 0 4px 15px rgba(0,0,0,0.5); background: #000; }
            .status { margin-top: 15px; font-size: 1.2rem; color: #4caf50; }
            .loading { color: #888; font-size: 0.9rem; }
        </style>
    </head>
    <body>
        <h1>🎥 实时人脸识别监控</h1>
        <div class="container">
            <!-- 关键：src 必须指向 /video_feed -->
            <img src="/video_feed" alt="Loading stream...">
        </div>
        <div class="status">系统运行中... 请等待 AI 连接</div>
        <div class="loading">如果画面一直不显示，请检查 ESP32 是否已上电并连接 Wi-Fi</div>
    </body>
    </html>
    """
    return html_content


if __name__ == '__main__':
    # 启动 AI 线程
    ai_thread = threading.Thread(target=ai_processing_loop, daemon=True)
    ai_thread.start()

    print(">>> 🚀 Flask 服务器启动中...")
    print(">>> 🌐 请在浏览器访问: http://127.0.0.1:5000")
    app.run(host='0.0.0.0', port=5000, threaded=True)