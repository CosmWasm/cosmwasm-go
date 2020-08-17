cmake_minimum_required(VERSION 3.16)
project(cw_compiler)

set(CMAKE_CXX_STANDARD 14)

add_executable(cw_compiler main.cpp BaseCheck.cpp BaseCheck.h error.h)