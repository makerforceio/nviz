## Assuming standard python main files
## File: main; Function: main.main(args); Args: main.args
import model3 as main

import os
import sys
import threading
import requests
import json
import uuid
import time

## Custom name for AI, URL for nvis
def wrapper(name, url, train_args = None):
    if url[len(url)-1] != '/':
        url += '/'
    ai_id = uuid.uuid4()
    name = os.path.splitext(os.path.basename(name))[0]

    print("Name: {}".format(name))
    print("UUID: {}".format(ai_id))
    print("URL: {}".format(url))

    train_thread = threading.Thread(target=main.main, args=(train_args,))
    train_thread.start()

    while not main.init_done:
        0

    args = {
        'name' : name,
        'args' : vars(main.args)
        }
        
    requests.put(url + "api/ai/{}".format(ai_id), data=json.dumps(args))
    
    render_thread = threading.Thread(target=render, args=(url, ai_id,))
    render_thread.start()

def render(url, ai_id):
    while True:
        progress = {
            'time'          : time.time(),
            'epoch'         : main.epoch,
            'training_loss' : main.training_loss,
            'stats'         : main.stats,
    }
        out = None
        if main.render:
            out = main.render()
                
        requests.post(url + "api/ai/{}/update".format(ai_id), data=json.dumps(progress))
        if out is not None:
            requests.post(url + "api/ai/{}/update/image".format(ai_id), data=out)
            
if __name__ == "__main__":
    wrapper(os.getenv('AI_NAME', sys.argv[0]), os.getenv('URL', 'localhost:8080'), sys.argv)
