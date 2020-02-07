#include "myclass.h"
#include <iostream>
#include <appdynamics.h>
#include <string>

extern "C" void add_to_call_graph(uintptr_t in, char *value)
{
    appd_bt_handle bt = (appd_bt_handle)in;

    appd::sdk::BT bt2(bt);
	
    appd::sdk::CallGraph callGraph(bt2, "Class1", "main", value, 276, 100, APPD_FRAME_TYPE_CPP);
    callGraph.add_to_snapshot();

    appd_bt_end(bt);

}
