## Assuming standard python main files
## File: main; Function: main.main(args); Args: main.args
import main

import os
import sys
import threading
import requests
import uuid
import time

## Custom name for AI, URL for nvis
def wrapper(name, url, train_args = None):
    if url[len(url)-1] != '/':
        url += '/'
    ai_id = uuid.uuid4()
    name = os.path.splitext(os.path.basename(name))

    print("Name: {}".format(name))
    print("UUID: {}".format(ai_id))
    print("URL: {}".format(url))

    train_thread = threading.Thread(target=main.main, args=(train_args))
    train_thread.start()

    requests.put("URL ARGS", data=main.args)

    render_thread = threading.Thread(target=render)
    render_thread.start()

def render():
    while True:
        if main.model is not None:
            progress = {
                    'time'          : time.time(),
                    'epoch'         : main.epoch,
                    'training_loss' : main.training_loss,
                    }
            out = main.render()
            
            requests.post("URL Data", data=progress)
            requests.post("URL Image", data=out)
            

if __name__ == "__main__":
    wrapper(os.getenv('AI_NAME', sys.argv[0]), os.getenv('URL', 'localhost:8080'), sys.argv)
