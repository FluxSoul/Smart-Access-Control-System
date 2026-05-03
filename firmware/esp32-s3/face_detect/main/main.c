#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "nvs_flash.h"
#include "led.h"
#include "my_spi.h"
#include "myiic.h"
#include "spilcd.h"
#include "esp_log.h"
#include "lwip/sockets.h"
#include "camera_config.h"
#include "wifi_config.h"
#include "camera_streamer.h"
#include "esp_heap_caps.h"
#include <stdio.h>


/**
 * @brief       程序入口
 * @param       无
 * @retval      无
 */
void app_main(void)
{
    esp_err_t ret;

    ret = nvs_flash_init();     /* 初始化NVS */
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ESP_ERROR_CHECK(nvs_flash_init());
    }

    led_init();                 /* LED初始化 */
    my_spi_init();              /* SPI初始化 */ 
    myiic_init();               /* MYIIC初始化 */
    xl9555_init();              /* XL9555初始化 */
    spilcd_init();              /* SPILCD初始化 */
    init_camera();              /* 初始化摄像头 */

    vTaskDelay(pdMS_TO_TICKS(2000));

    wifi_sta_init();

    ESP_LOGI("main", "wifi success");

    start_camera_stream_server();

    ESP_LOGI("main", "server run success");
    spilcd_show_string(0, 0, 60, 16, 16, "server run suceess", BLUE);

    while(1)
    {   
        // ESP_LOGI(TAG, "main check");
        vTaskDelay(100 / portTICK_PERIOD_MS);
    }
}
