#include <errno.h>
#include <modbus/modbus.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "MySqlWrite.h"

#include "datastruct.h"

extern bool silent;

int ModbusRtuABB(bool updatecounter) {

  uint16_t tab_reg[66];
  int rc;
  modbus_t *ctx = NULL;

  struct instValues data;

  ctx = modbus_new_rtu("/dev/ttyUSB0", 19200, 'N', 8, 1);

  if (ctx == NULL) {
    fprintf(stderr, "Unable to allocate libmodbus context\n");
    return -1;
  }

  int slaveAddr = 1;

  modbus_set_debug(ctx, FALSE);

  modbus_set_slave(ctx, slaveAddr);

  // modbus_set_error_recovery(ctx,
  //		MODBUS_ERROR_RECOVERY_LINK | MODBUS_ERROR_RECOVERY_PROTOCOL);

  if (modbus_connect(ctx) == -1) {
    fprintf(stderr, "Connection failed: %s\n", modbus_strerror(errno));
    modbus_free(ctx);
    return -1;
  }

  rc = modbus_read_registers(ctx, 0x5B00, sizeof(tab_reg) / sizeof(tab_reg[0]),
                             tab_reg);
  if (rc == -1) {
    fprintf(stderr, "%s\n", modbus_strerror(errno));
    modbus_close(ctx);
    modbus_free(ctx);
    return -1;
  } else {

    if (!silent) {
      for (int j = 0; j < rc; j++) {
        printf("%04X\t: %04X\n", 0x5B00 + j, tab_reg[j]);
      }
    }

    // Map to struct
    data.VoltL1_N = (tab_reg[0] << 16 | tab_reg[1]);
    data.VoltL2_N = (tab_reg[2] << 16 | tab_reg[3]);
    data.VoltL3_N = (tab_reg[4] << 16 | tab_reg[5]);

    data.VoltL1_L2 = (tab_reg[6] << 16 | tab_reg[7]);
    data.VoltL3_L2 = (tab_reg[8] << 16 | tab_reg[9]);
    data.VoltL1_L3 = (tab_reg[10] << 16 | tab_reg[11]);

    data.CurrentL1 = (tab_reg[12] << 16 | tab_reg[13]);
    data.CurrentL2 = (tab_reg[14] << 16 | tab_reg[15]);
    data.CurrentL3 = (tab_reg[16] << 16 | tab_reg[17]);
    data.CurrentN = (tab_reg[18] << 16 | tab_reg[19]);

    data.ActivePowerTotal = (tab_reg[20] << 16 | tab_reg[21]);
    data.ActivePowerL1 = (tab_reg[22] << 16 | tab_reg[23]);
    data.ActivePowerL2 = (tab_reg[24] << 16 | tab_reg[25]);
    data.ActivePowerL3 = (tab_reg[26] << 16 | tab_reg[27]);

    data.ReactivePowerTotal = (tab_reg[28] << 16 | tab_reg[29]);
    data.ReactivePowerL1 = (tab_reg[30] << 16 | tab_reg[31]);
    data.ReactivePowerL2 = (tab_reg[32] << 16 | tab_reg[33]);
    data.ReactivePowerL3 = (tab_reg[34] << 16 | tab_reg[35]);

    data.ApparentPowerTotal = (tab_reg[36] << 16 | tab_reg[37]);
    data.ApparentPowerL1 = (tab_reg[38] << 16 | tab_reg[39]);
    data.ApparentPowerL2 = (tab_reg[40] << 16 | tab_reg[41]);
    data.ApparentPowerL3 = (tab_reg[42] << 16 | tab_reg[43]);

    data.Frequency = tab_reg[44];

    data.PhaseAnglePowerTotal = tab_reg[45];
    data.PhaseAnglePowerL1 = tab_reg[46];
    data.PhaseAnglePowerL2 = tab_reg[47];
    data.PhaseAnglePowerL3 = tab_reg[48];

    data.PhaseAngleVoltageL1 = tab_reg[49];
    data.PhaseAngleVoltageL2 = tab_reg[50];
    data.PhaseAngleVoltageL3 = tab_reg[51];

    data.PhaseAngleCurrentL1 = tab_reg[55];
    data.PhaseAngleCurrentL2 = tab_reg[56];
    data.PhaseAngleCurrentL3 = tab_reg[57];

    data.PowerFactorTotal = tab_reg[58];
    data.PowerFactorL1 = tab_reg[59];
    data.PowerFactorL2 = tab_reg[60];
    data.PowerFactorL3 = tab_reg[61];

    data.CurrentQuadrantTotal = tab_reg[62];
    data.CurrentQuadrantL1 = tab_reg[63];
    data.CurrentQuadrantL2 = tab_reg[64];
    data.CurrentQuadrantL3 = tab_reg[65];

    if (!silent) {
      printf("L1-N: %.1f Volt\n", (double)data.VoltL1_N / 10);
      printf("L2-N: %.1f Volt\n", (double)data.VoltL2_N / 10);
      printf("L3-N: %.1f Volt\n", (double)data.VoltL3_N / 10);

      printf("L1-L2: %.1f Volt\n", (double)data.VoltL1_L2 / 10);
      printf("L3-L2: %.1f Volt\n", (double)data.VoltL3_L2 / 10);
      printf("L1-L3: %.1f Volt\n", (double)data.VoltL1_L3 / 10);

      printf("L1: %.2f Ampere\n", (double)data.CurrentL1 / 100);
      printf("L2: %.2f Ampere\n", (double)data.CurrentL2 / 100);
      printf("L3: %.2f Ampere\n", (double)data.CurrentL3 / 100);
      // Current N 18+19 not avail.

      printf("Active power Total : %.2f Watt\n",
             (double)data.ActivePowerTotal / 100);
      printf("Active power L1: %.2f Watt\n", (double)data.ActivePowerL1 / 100);
      printf("Active power L2: %.2f Watt\n", (double)data.ActivePowerL2 / 100);
      printf("Active power L3: %.2f Watt\n", (double)data.ActivePowerL3 / 100);

      printf("Frequency: %.2f Hz\n", (double)data.Frequency / 100);

      printf("Power Factor Total: %.3f pF\n",
             (double)data.PowerFactorTotal / 1000);
    }

    // SaveLiveData(data);
  }

  if (updatecounter) {

    rc = modbus_read_registers(ctx, 0x5000, 4, tab_reg);
    if (rc == -1) {
      fprintf(stderr, "%s\n", modbus_strerror(errno));
      modbus_close(ctx);
      modbus_free(ctx);
      return -1;
    } else {

      double imported = (double)(tab_reg[0] << 16 | tab_reg[1]) / 100;
      printf("Imported %.2f kWh\n", imported);

      double exported = (double)(tab_reg[2] << 16 | tab_reg[3]) / 100;
      printf("Exported %.2f kWh\n", exported);

      // SaveCounter((int)exported, (int)imported);
    }
  }

  // Registers 0x5170 - 32 have nothing 0xFFFF;
  // Registers 0x5460 have 6 32bit values

  /*
   rc = modbus_read_registers(ctx,0x5B00,64,tab_reg);
   if (rc == -1) {
   fprintf(stderr, "%s\n", modbus_strerror(errno));
   modbus_close(ctx);
   modbus_free(ctx);
   return -1;
   } else {

   printf("START: 0x5B00\n");
   for(j=0;j<rc;j++){
   printf("%04X\t: %04X\n",0x5B00+j,tab_reg[j]);
   }





   }
   */

  modbus_close(ctx);
  modbus_free(ctx);

  return EXIT_SUCCESS;
}
