#!/bin/bash

fileloc=$(mage build)

if [ $? == 0 ]
then
    echo "Running" $fileloc
    chmod +x $fileloc
    {
        cd ./run
        eval "../$fileloc"
    }
fi


