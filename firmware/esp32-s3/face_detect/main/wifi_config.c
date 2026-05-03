#include "wifi_config.h"

EventGroupHandle_t   wifi_event;
char lcd_buff[100] = {0};

/**
 * @brief       链接显示
 * @param       flag:2->链接;1->链接失败;0->再链接中
 * @retval      无
 */
void connet_display(uint8_t flag)
{
    if(flag == 2)
    {
        spilcd_fill(0,90,320,240,WHITE);
        sprintf(lcd_buff, "ssid:%s",DEFAULT_SSID);
        spilcd_show_string(0, 90, 240, 16, 16, lcd_buff, BLUE);
        sprintf(lcd_buff, "psw:%s",DEFAULT_PWD);
        spilcd_show_string(0, 110, 240, 16, 16, lcd_buff, BLUE);
    }
    else if (flag == 1)
    {
        spilcd_show_string(0, 90, 240, 16, 16, "wifi connecting fail", BLUE);
    }
    else
    {
        spilcd_show_string(0, 90, 240, 16, 16, "wifi connecting......", BLUE);
    }
}

/**
 * @brief       WIFI链接糊掉函数
 * @param       arg:传入网卡控制块
 * @param       event_base:WIFI事件
 * @param       event_id:事件ID
 * @param       event_data:事件数据
 * @retval      无
 */
void wifi_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
{
    static int s_retry_num = 0;

    /* 扫描到要连接的WIFI事件 */
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START)
    {
        connet_display(0);
        esp_wifi_connect();
    }
    /* 连接WIFI事件 */
    else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_CONNECTED)
    {
        connet_display(2);
    }
    /* 连接WIFI失败事件 */
    else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED)
    {
        /* 尝试连接 */
        if (s_retry_num < 20)
        {
            esp_wifi_connect();
            s_retry_num++;
            ESP_LOGI(TAG, "retry to connect to the AP");
        }
        else
        {
            xEventGroupSetBits(wifi_event, WIFI_FAIL_BIT);
        }

        ESP_LOGI(TAG,"connect to the AP fail");
    }
    /* 工作站从连接的AP获得IP */
    else if(event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP)
    {
        ip_event_got_ip_t* event = (ip_event_got_ip_t*) event_data;
        ESP_LOGI(TAG, "static ip:" IPSTR, IP2STR(&event->ip_info.ip));
        s_retry_num = 0;
        xEventGroupSetBits(wifi_event, WIFI_CONNECTED_BIT);
    }
}

/**
 * @brief       WIFI初始化
 * @param       无
 * @retval      无
 */
void wifi_sta_init(void)
{
    static esp_netif_t *sta_netif = NULL;
    wifi_event= xEventGroupCreate();    /* 创建一个事件标志组 */
    /* 网卡初始化 */
    ESP_ERROR_CHECK(esp_netif_init());
    /* 创建新的事件循环 */
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    sta_netif= esp_netif_create_default_wifi_sta();
    assert(sta_netif);
    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK( esp_event_handler_register(WIFI_EVENT, ESP_EVENT_ANY_ID, &wifi_event_handler, NULL) );
    ESP_ERROR_CHECK( esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &wifi_event_handler, NULL) );
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));    
    wifi_config_t  wifi_config = WIFICONFIG();
    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK( esp_wifi_set_config(ESP_IF_WIFI_STA, &wifi_config) );
    ESP_ERROR_CHECK(esp_wifi_start());

    /* 等待链接成功后、ip生成 */
    EventBits_t bits = xEventGroupWaitBits(wifi_event,
            WIFI_CONNECTED_BIT | WIFI_FAIL_BIT,
            pdFALSE,
            pdFALSE,
            portMAX_DELAY);

    /* 判断连接事件 */
    if (bits & WIFI_CONNECTED_BIT)
    {
        ESP_LOGI(TAG, "connected to ap SSID:%s password:%s",
                 DEFAULT_SSID, DEFAULT_PWD);
    }
    else if (bits & WIFI_FAIL_BIT)
    {
        connet_display(1);
        ESP_LOGI(TAG, "Failed to connect to SSID:%s, password:%s",
                 DEFAULT_SSID, DEFAULT_PWD);
    }
    else
    {
        ESP_LOGE(TAG, "UNEXPECTED EVENT");
    }
}