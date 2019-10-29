#include <errno.h>
#include <mysql/mysql.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "datastruct.h"

MYSQL *conn;
MYSQL_RES *res;
MYSQL_ROW row;

char mysql_host[100] = "192.168.?.?";
char mysql_user[32] = "abbmeter";
char mysql_pass[32] = "?";
char mysql_db[100] = "?";
unsigned int mysql_port = 3306;

int SaveLiveData(struct instValues data) {

  conn = mysql_init(NULL);
  if (conn == NULL) {
    perror("mysql_init");
  } else {

    if (!mysql_real_connect(conn, mysql_host, mysql_user, mysql_pass, mysql_db,
                            mysql_port, NULL, 0)) {
      fprintf(stderr, "%s\n", mysql_error(conn));
    } else {
      char stmt_buf[1024];
      sprintf(stmt_buf,
              "INSERT INTO `energiemeter_abb` (`moment`,"
              "`voltL1`, `voltL2`, `voltL3`,"
              "`voltL1L2`, `voltL3L2`, `voltL1L3`,"
              "`currentL1`, `currentL2`, `currentL3`,"
              "`activePowerTotal`, `activePowerL1`, `activePowerL2`, "
              "`activePowerL3`,"
              "`frequency`, `powerFactorTotal`) VALUES ( NOW(3), "
              " %d, %d, %d, "
              " %d, %d, %d, "
              " %d, %d, %d, "
              " %d, %d, %d, %d, "
              " %d, %d);",
              data.VoltL1_N, data.VoltL2_N, data.VoltL3_N, data.VoltL1_L2,
              data.VoltL3_L2, data.VoltL1_L3, data.CurrentL1, data.CurrentL2,
              data.CurrentL3, data.ActivePowerTotal, data.ActivePowerL1,
              data.ActivePowerL2, data.ActivePowerL3, data.Frequency,
              data.PowerFactorTotal);

      int nQueryResult = mysql_query(conn, stmt_buf);
      if (nQueryResult > 0) {
        fprintf(stdout, "%s\n", mysql_error(conn));
      } else {
        printf("Rows added %d\n", nQueryResult);
      }

      // res = mysql_use_result(conn);

      // while((row = mysql_fetch_row(res)) != NULL){
      //	printf("Database: %s\n",row[0]);
      //	//int dumpresult = system("")
      //}

      // mysql_free_result(res);
      mysql_close(conn);
    }
  }

  return EXIT_SUCCESS;
}

int SaveCounter(int counter1, int counter2) {
  conn = mysql_init(NULL);
  if (conn == NULL) {
    perror("mysql_init");
  } else {

    if (!mysql_real_connect(conn, mysql_host, mysql_user, mysql_pass, mysql_db,
                            mysql_port, NULL, 0)) {
      fprintf(stderr, "%s\n", mysql_error(conn));
    } else {
      char stmt_buf[1024];
      sprintf(stmt_buf,
              "INSERT INTO `meterstanden` (`meter`, `moment`, `stand1`, "
              "`stand2`) VALUES ('ABB', NOW(), %d, %d);",
              counter1, counter2);

      int nQueryResult = mysql_query(conn, stmt_buf);
      if (nQueryResult > 0) {
        fprintf(stdout, "%s\n", mysql_error(conn));
      } else {
        printf("Rows added %d\n", nQueryResult);
      }

      // res = mysql_use_result(conn);

      // while((row = mysql_fetch_row(res)) != NULL){
      //	printf("Database: %s\n",row[0]);
      //	//int dumpresult = system("")
      //}

      // mysql_free_result(res);
      mysql_close(conn);
    }
  }

  return EXIT_SUCCESS;
}
