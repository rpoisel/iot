add_executable(AbbB23Energymeter
  main.c
  ModbusRtuClient.c
)

target_include_directories(AbbB23Energymeter PUBLIC
  ${CMAKE_CURRENT_LIST_DIR}
)

target_link_libraries(AbbB23Energymeter PUBLIC
  modbus
)

target_compile_options(AbbB23Energymeter PUBLIC
  -Wall -Wextra
)
