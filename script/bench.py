import os

lines_count = 250000
log_dir = './bench_logs'
log_file_name = log_dir + '/test.log'
config_file = "config.ini"

logshot_binary = 'logshot'
binary = "go run main.go"
msg = "test string one\n"

run_params = "-config=%s -readall -logfile=%s &>/dev/null" % (config_file, log_file_name)


def bench(logs_count=1):
    os.system("rm -rf %s; mkdir -p %s" % (log_dir, log_dir))
    for x in range(0, logs_count):
        with open(log_file_name + str(x), "a") as myfile:
            myfile.write(msg * lines_count)

    os.system("/bin/bash -c '%s %s'" % (binary, run_params))
    os.system("echo '\n'")
    os.system("rm -f %s" % (log_file_name + '*'))


if __name__ == '__main__':
    print("with 1 file containing %s matching lines each" % lines_count)
    bench()
