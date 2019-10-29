/*
 ============================================================================
 Name        : WagoClient.c
 Author      : Robin Jansman
 Version     :
 Copyright   :
 Description :
 ============================================================================
 */

#include <errno.h>
#include <modbus/modbus.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "ModbusRtuClient.h"

int startswith(const char *prefix, const char *string);

bool silent = false;

const int cCounter = 60 * 60 * 6; // every 6 hours

int main(int argc, char *argv[]) {

  int i = 0;
  for (i = 0; i < argc; i++) {
    if (startswith("--silent", argv[i])) {
      silent = true;
    }
  }

  //	int miliseconds = 1000;
  //	int intervalus = miliseconds*1000;

  int counter = 0;

  while (true) {

    ModbusRtuABB(counter == 0);

    if (counter <= 0) {
      counter = cCounter;
    }
    counter = counter - 1;

    // usleep(intervalus);
    sleep(1);
  }
}

int startswith(const char *prefix, const char *string) {
  size_t lenstring = strlen(string);
  size_t lenprefix = strlen(prefix);
  if (lenstring < lenprefix) {
    return 0;
  } else {
    return (strncmp(string, prefix, lenprefix) == 0);
  }
}
