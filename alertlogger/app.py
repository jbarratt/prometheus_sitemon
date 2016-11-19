from flask import Flask, request
import os
app = Flask(__name__)


@app.route('/', methods=['POST'])
def root_handler():
    if 'LOG_ALERT_PATH' in os.environ:
        print(str(request.data))
        with open(os.environ['LOG_ALERT_PATH'], "a") as log:
            log.write(str(request.data))
            log.write("\n")
    else:  # just go to stdou
        print(str(request.data))
    return "{\"status\": \"ok\"}\n"
