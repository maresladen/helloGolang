#!/bin/sh
useport=5000
testData=""
for PID in $(lsof -i:$useport |awk '{print $2}'); do
if [ $PID != "PID" ]; then
    if [ -z $PID ]; then
        break;
        else
        testData=$PID
        kill $PID
    fi
fi
done
cd /Users/BetaFun/CodeWork/DotNet/asptest
dotnet run