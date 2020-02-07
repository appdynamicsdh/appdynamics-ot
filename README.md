# appdynamics-ot

export APPD_SDK_HOME=/home/djhope/appdynamics-cpp-sdk

export CGO_CFLAGS="-I $APPD_SDK_HOME/include"

export CGO_LDFLAGS="-L $APPD_SDK_HOME/lib -l appdynamics -Wl,--no-as-needed -ldl -lmyclass"

export LD_LIBRARY_PATH=$APPD_SDK_HOME/lib
