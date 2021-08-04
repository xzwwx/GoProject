#!/bin/sh

startwork(){
  nohup $PWD/bin/rcenterserver -config=$PWD/bin/Config/Config.json &
  sleep 2s
  nohup $PWD/bin/loginserver -config=$PWD/bin/Config/Config.json &
  nohup $PWD/bin/roomserver -config=$PWD/bin/Config/Config.json &


}