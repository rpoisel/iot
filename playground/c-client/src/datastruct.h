#include <errno.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#ifndef INSTVALUES_H
#define INSTVALUES_H

struct instValues {

  uint32_t VoltL1_N; // 0x5B00
  uint32_t VoltL2_N; // 0x5B02
  uint32_t VoltL3_N; // 0x5B04

  uint32_t VoltL1_L2; // 0x5B06
  uint32_t VoltL3_L2; // 0x5B08
  uint32_t VoltL1_L3; // 0x5B0A

  uint32_t CurrentL1; // 0x5B0C
  uint32_t CurrentL2; // 0x5B0E
  uint32_t CurrentL3; // 0x5B10
  uint32_t CurrentN;  // 0x5B12, not available

  int32_t ActivePowerTotal; // 0x5B14
  int32_t ActivePowerL1;    // 0x5B16
  int32_t ActivePowerL2;    // 0x5B18
  int32_t ActivePowerL3;    // 0x5B1A

  int32_t ReactivePowerTotal; // 0x5B1C
  int32_t ReactivePowerL1;    // 0x5B1E
  int32_t ReactivePowerL2;    // 0x5B20
  int32_t ReactivePowerL3;    // 0x5B22

  int32_t ApparentPowerTotal; // 0x5B24
  int32_t ApparentPowerL1;    // 0x5B26
  int32_t ApparentPowerL2;    // 0x5B28
  int32_t ApparentPowerL3;    // 0x5B2A

  uint16_t Frequency; // 0x5B2C

  int16_t PhaseAnglePowerTotal; // 0x52B2D
  int16_t PhaseAnglePowerL1;    // 0x52B2E
  int16_t PhaseAnglePowerL2;    // 0x52B2F
  int16_t PhaseAnglePowerL3;    // 0x52B30

  int16_t PhaseAngleVoltageL1; // 0x5231
  int16_t PhaseAngleVoltageL2; // 0x5232
  int16_t PhaseAngleVoltageL3; // 0x5233

  // GAP 3 bytes

  int16_t PhaseAngleCurrentL1; // 0x5237
  int16_t PhaseAngleCurrentL2; // 0x5238
  int16_t PhaseAngleCurrentL3; // 0x5239

  int16_t PowerFactorTotal; // 0x523A		(wel)
  int16_t PowerFactorL1;    // 0x523B
  int16_t PowerFactorL2;    // 0x523C
  int16_t PowerFactorL3;    // 0x523D

  uint16_t CurrentQuadrantTotal; // 0x5B3E
  uint16_t CurrentQuadrantL1;    // 0x5B3F
  uint16_t CurrentQuadrantL2;    // 0x5B40
  uint16_t CurrentQuadrantL3;    // 0x5B41
};

#endif
