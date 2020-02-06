package appd

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <appdynamics.h>
#include <myclass.h>


// global configuration
struct appd_config *cfg;

void printpair(char *name, char *value) {
    printf("%s = %s\n", name, value);
}

uintptr_t bt_handle_to_int(appd_bt_handle bt_handle) {
    return (uintptr_t)bt_handle;
}
appd_bt_handle bt_int_to_handle(uintptr_t bt) {
    return (appd_bt_handle)bt;
}

uintptr_t exit_handle_to_int(appd_exitcall_handle exitcall_handle) {
    return (uintptr_t)exitcall_handle;
}
appd_exitcall_handle exit_int_to_handle(uintptr_t exit) {
    return (appd_exitcall_handle)exit;
}
*/
import "C"

import (
    "log"
    "unsafe"
    "net/http"
    "strconv"
    "runtime"
    "os/exec"
    "time"
    "strings"
)

/*
* Some constants from appdynamics.h
*/
const CORRELATION_HEADER_NAME = C.APPD_CORRELATION_HEADER_NAME

const BACKEND_HTTP = C.APPD_BACKEND_HTTP
const BACKEND_DB = C.APPD_BACKEND_DB
const BACKEND_CACHE = C.APPD_BACKEND_CACHE
const BACKEND_RABBITMQ = C.APPD_BACKEND_RABBITMQ
const BACKEND_WEBSERVICE = C.APPD_BACKEND_WEBSERVICE

const ERROR_LEVEL_NOTICE = C.APPD_LEVEL_NOTICE
const ERROR_LEVEL_WARNING = C.APPD_LEVEL_WARNING
const ERROR_LEVEL_ERROR = C.APPD_LEVEL_ERROR

const APPD_BT = "APPD_BT"

type ID_properties_map map[string]string

// do not free the C.CString as the config struct is only pointers
func Init(appName string, controllerKey string) {
    log.Println("Init called")


    C.cfg = C.appd_config_init() // appd_config_init() resets the configuration object and pass back an handle/pointer
    C.appd_config_set_app_name(C.cfg, C.CString(appName))
    C.appd_config_set_controller_access_key(C.cfg, C.CString(controllerKey))




}

func SetTierName(name string) {
    //C.cfg.tier_name =
    C.appd_config_set_tier_name(C.cfg, C.CString(name))

}

func SetNodeName(name string) {
    //C.cfg.node_name = C.CString(name)

    C.appd_config_set_node_name(C.cfg, C.CString(name))
}

func SetControllerHost(host string) {
    //C.cfg.controller.host = C.CString(host)

    C.appd_config_set_controller_host(C.cfg, C.CString(host))
}

func SetControllerPort(port int16) {
    //C.cfg.controller.port =


    C.appd_config_set_controller_port(C.cfg, C.ushort(port))
}

func SetControllerAccount(account string) {
    C.appd_config_set_controller_account(C.cfg, C.CString(account))
}

func SetControllerUseSSL(ssl int16) {
    //C.cfg.controller.use_ssl =

    C.appd_config_set_controller_use_ssl(C.cfg, C.uint(ssl))
}

// Proxy
func SetControllerProxyHost(host string) {
    //C.cfg.controller.http_proxy.host = C.CString(host)
}

func SetControllerProxyPort(port int16) {
    //C.cfg.controller.http_proxy.port = C.ushort(port)
}

func SetControllerProxyUsername(username string) {
    //C.cfg.controller.http_proxy.username = C.CString(username)
}

func SetControllerProxyPasswordFile(password_file string) {
    //C.cfg.controller.http_proxy.password_file = C.CString(password_file)
}

// Agent
func SetAgentProxyControlPort(port int16) {
    //C.cfg.agent_proxy.tcp_control_port = C.ushort(port)
}

func SetAgentProxyRequestPort(port int16) {
    //C.cfg.agent_proxy.tcp_request_port = C.ushort(port)
}

func SetAgentProxyReportingPort(port int16) {
    //C.cfg.agent_proxy.tcp_reporting_port = C.ushort(port)
}

func SetAgentProxyCommDir(dir string) {
    //C.cfg.agent_proxy.ipc_comm_dir = C.CString(dir)
}

func SetInitTimeout(timeout int) {
    //C.cfg.init_timeout_ms = C.int(timeout)
}

func Sdk_init() int {
    //C.dump_config(&C.cfg)
    rc := C.appd_sdk_init(C.cfg)
    return int(rc)
}

func Sdk_term() {
    C.appd_sdk_term()
}

/*
* BT
*/
func BT_begin(name string, correlation string) uint64 {
    name_c := C.CString(name)
    correlation_c := C.CString(correlation)
    defer C.free(unsafe.Pointer(name_c))
    defer C.free(unsafe.Pointer(correlation_c))

    bt := C.appd_bt_begin(name_c, correlation_c)

    return uint64(C.bt_handle_to_int(bt))
}

func BT_end(bt uint64) {
    callGraph(bt)
    C.appd_bt_end(C.bt_int_to_handle(C.uintptr_t(bt)))
}

func BT_override_start_time_ms(bt uint64, start uint64) {
	C.appd_bt_override_start_time_ms(C.bt_int_to_handle(C.uintptr_t(bt)), C.long(start))
}

func BT_override_time_ms(bt uint64, timeMS uint64) {
	C.appd_bt_override_time_ms(C.bt_int_to_handle(C.uintptr_t(bt)), C.long(timeMS))
}

func BT_set_url(bt uint64, name string) {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    C.appd_bt_set_url(C.bt_int_to_handle(C.uintptr_t(bt)), name_c)
}

func BT_is_snapshotting(bt uint64) int {
    result := C.appd_bt_is_snapshotting(C.bt_int_to_handle(C.uintptr_t(bt)))
    return int(result)
}

func BT_add_user_data(bt uint64, key string, value string) {
    key_c := C.CString(key)
    value_c := C.CString(value)
    defer C.free(unsafe.Pointer(key_c))
    defer C.free(unsafe.Pointer(value_c))

    C.appd_bt_add_user_data(C.bt_int_to_handle(C.uintptr_t(bt)), key_c, value_c)
}

func BT_add_error(bt uint64, level uint32, message string, mark_bt_as_error int) {
    message_c := C.CString(message)
    defer C.free(unsafe.Pointer(message_c))

    C.appd_bt_add_error(C.bt_int_to_handle(C.uintptr_t(bt)), level, message_c, C.int(mark_bt_as_error))
}

func BT_store(bt uint64, guid string) {
    guid_c := C.CString(guid)
    defer C.free(unsafe.Pointer(guid_c))

    C.appd_bt_store(C.bt_int_to_handle(C.uintptr_t(bt)), guid_c)
}

func BT_get(guid string) uint64 {
    guid_c := C.CString(guid)
    defer C.free(unsafe.Pointer(guid_c))

    bt := C.appd_bt_get(guid_c)

    return uint64(C.bt_handle_to_int(bt))
}

/*
* Backend
*/
func Backend_declare(betype string, name string) {
    betype_c := C.CString(betype)
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(betype_c))
    defer C.free(unsafe.Pointer(name_c))

    log.Println("Backend decalre", betype, name)
    C.appd_backend_declare(C.CString(betype), C.CString(name))
}

func Backend_set_identifying_property(name string, key string, value string) int {
    name_c := C.CString(name)
    key_c := C.CString(key)
    value_c := C.CString(value)
    defer C.free(unsafe.Pointer(name_c))
    defer C.free(unsafe.Pointer(key_c))
    defer C.free(unsafe.Pointer(value_c))

    rc := C.appd_backend_set_identifying_property(C.CString(name), C.CString(key), C.CString(value))
    return int(rc)
}

func Backend_set_identifying_properties(name string, props ID_properties_map) int {
    var rc C.int
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    for key, value := range props {
        key_c := C.CString(key)
        value_c := C.CString(value)
        rc = C.appd_backend_set_identifying_property(name_c, key_c, value_c)
        C.free(unsafe.Pointer(key_c))
        C.free(unsafe.Pointer(value_c))
        if(rc != 0) {
            break
        }
    }
    return int(rc)
}

func Backend_add(name string) int {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    rc := C.appd_backend_add(name_c)
    return int(rc)
}

func Backend_prevent_agent_resolution(name string) int {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    rc := C.appd_backend_prevent_agent_resolution(name_c)
    return int(rc)
}

/*
* Exit
*/
func Exitcall_begin(bt uint64, name string) uint64 {
    name_c := C.CString(name)
    defer C.free(unsafe.Pointer(name_c))

    bt_h := C.bt_int_to_handle(C.uintptr_t(bt))
    exit := C.appd_exitcall_begin(bt_h, name_c)

    return uint64(C.exit_handle_to_int(exit))
}

func Exitcall_end(exit uint64) {
    exit_h := C.exit_int_to_handle(C.uintptr_t(exit))
    C.appd_exitcall_end(exit_h)
}

func Exitcall_set_details(exit uint64, details string) int {
    details_c := C.CString(details)
    defer C.free(unsafe.Pointer(details_c))

    rc := C.appd_exitcall_set_details(C.exit_int_to_handle(C.uintptr_t(exit)), details_c)
    return int(rc)
}

func Exitcall_add_error(exit uint64, error_level uint32, message string, mark_bt_as_error int) {
    message_c := C.CString(message)
    defer C.free(unsafe.Pointer(message_c))

    C.appd_exitcall_add_error(C.exit_int_to_handle(C.uintptr_t(exit)), error_level, message_c, C.int(mark_bt_as_error))
}

func Exitcall_get_correlation_header(exit uint64) string {
    header := C.appd_exitcall_get_correlation_header(C.exit_int_to_handle(C.uintptr_t(exit)))
    return C.GoString(header)
}

func Exitcall_store(exit uint64, guid string) {
    guid_c := C.CString(guid)
    defer C.free(unsafe.Pointer(guid_c))

    C.appd_exitcall_store(C.exit_int_to_handle(C.uintptr_t(exit)), guid_c)
}

func Exitcall_get(guid string) uint64 {
    guid_c := C.CString(guid)
    defer C.free(unsafe.Pointer(guid_c))

    exit := C.appd_exitcall_get(guid_c)

    return uint64(C.exit_handle_to_int(exit))
}

/*
* HTTP wrappers
*/
/*
* Wrappers for http handlers
* name is the name used for the BT
*/
func WrapHandle(name string, pattern string, handler http.Handler) (string, http.Handler) {
    return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // get correlation id from header
        appd_correlation := r.Header.Get(CORRELATION_HEADER_NAME)
        bt := BT_begin(name, appd_correlation)
        defer BT_end(bt)
        BT_set_url(bt, r.URL.String())
        r.Header.Add(APPD_BT, strconv.FormatUint(bt, 10))
        handler.ServeHTTP(w, r)
    })
}

func WrapHandleFunc(name string, pattern string, handler func(http.ResponseWriter, *http.Request)) (string, func(http.ResponseWriter, *http.Request)) {
    p, h := WrapHandle(name, pattern, http.HandlerFunc(handler))
    return p, func(w http.ResponseWriter, r *http.Request) {
        h.ServeHTTP(w, r)
    }
}
var x string

func callGraph(bt uint64) {

    threads := runtime.GOMAXPROCS(0)

    for i := 0; i < threads; i++ {
        go func() {
            for {
                out, err := exec.Command("ls","-lah").Output()

                if err != nil {
                        log.Fatalf("cmd.Run() failed with %s\n", err)
                }
                x=string(out)
                time.Sleep(1000 * time.Millisecond)
            }
        }()
    }

        var david = strings.Split(x, "\n")
	for line := range david{
		C.add_to_call_graph(C.uintptr_t(bt),C.CString(string(line)));
	}


}

