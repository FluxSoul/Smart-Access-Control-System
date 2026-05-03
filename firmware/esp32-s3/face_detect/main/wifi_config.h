#ifndef WIFI_CONFIG_H
#define WIFI_CONFIG_H

#include "esp_wifi.h"
#include "esp_event.h"
#include "esp_mac.h"
#include "spilcd.h"
#include "esp_log.h"

/* 链接wifi名称 */
#define DEFAULT_SSID        "iPhoneYang"
/* wifi密码 */
#define DEFAULT_PWD         "wangyuting"

extern EventGroupHandle_t   wifi_event;
static const char *TAG = "static_ip";
extern char lcd_buff[100];

#define WIFI_CONNECTED_BIT  BIT0
#define WIFI_FAIL_BIT       BIT1

/* WIFI默认配置 */
#define WIFICONFIG()   {                            \
    .sta = {                                        \
        .ssid = DEFAULT_SSID,                       \
        .password = DEFAULT_PWD,                    \
        .threshold.authmode = WIFI_AUTH_WPA2_PSK,   \
    },                                              \
}

void connet_display(uint8_t flag);
void wifi_event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data);
void wifi_sta_init(void);

#endif