/*
 * MySqlWrite.h
 *
 *  Created on: Jul 7, 2018
 *      Author: robin
 */

#include "datastruct.h"

#ifndef MYSQLWRITE_H_
#define MYSQLWRITE_H_


int SaveLiveData(struct instValues data);

int SaveCounter(int counter1, int counter2);

#endif /* MYSQLWRITE_H_ */
