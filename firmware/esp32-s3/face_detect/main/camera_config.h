#ifndef CAMERA_CONFIG_H
#define CAMERA_CONFIG_H

#include "esp_camera.h"
#include "esp_log.h"
#include "xl9555.h"

/* 引脚配置 */
#define CAM_PIN_PWDN    GPIO_NUM_NC
#define CAM_PIN_RESET   GPIO_NUM_NC
#define CAM_PIN_VSYNC   GPIO_NUM_47
#define CAM_PIN_HREF    GPIO_NUM_48
#define CAM_PIN_PCLK    GPIO_NUM_45
#define CAM_PIN_XCLK    GPIO_NUM_NC
#define CAM_PIN_SIOD    GPIO_NUM_NC // 已经被定义
#define CAM_PIN_SIOC    GPIO_NUM_NC
#define CAM_PIN_D0      GPIO_NUM_4
#define CAM_PIN_D1      GPIO_NUM_5
#define CAM_PIN_D2      GPIO_NUM_6
#define CAM_PIN_D3      GPIO_NUM_7
#define CAM_PIN_D4      GPIO_NUM_15
#define CAM_PIN_D5      GPIO_NUM_16
#define CAM_PIN_D6      GPIO_NUM_17
#define CAM_PIN_D7      GPIO_NUM_18


#define CAM_PWDN(x)         do{ x ? \
                                (xl9555_pin_write(OV_PWDN_IO, 1)):       \
                                (xl9555_pin_write(OV_PWDN_IO, 0));       \
                            }while(0)

#define CAM_RST(x)          do{ x ? \
                                (xl9555_pin_write(OV_RESET_IO, 1)):       \
                                (xl9555_pin_write(OV_RESET_IO, 0));       \
                            }while(0)


extern camera_config_t camera_config;
esp_err_t init_camera(void);

#endif
