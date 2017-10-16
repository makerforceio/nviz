import argparse
import random
import io

import torch

import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim

from torch.autograd import Variable

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

from model import *

parser = argparse.ArgumentParser()
parser.add_argument('power', type=int, metavar='power')

init_done = False
model, args, epoch, training_loss, stats = None, None, None, None, {}
target_model = None
def main(argv):
    global init_done
    global model, args, epoch, training_loss
    global target_model

    args = parser.parse_args(args)

    model = Model(args.power)
    target_model = Model(args.power)

    optimizer = optim.Adam(model.parameters())

    init_done = True

    epoch = 0
    while True:
        optimizer.zero_grad()
        x = Variable(torch.Tensor([random.random()]))
        output = model(x)
        target = target_model(x)
        loss = F.mse_loss(output, target)
        loss.backward()
        loss = loss.data[0]
        optimizer.step()
        
        epoch += 1
        training_loss = loss

        print("\rEpoch: {}".format(epoch), end="")

targets = None
def render():
    global model, targets
    if model and target_model:
        model_clone = Model(model.polySize)
        model_clone.load_state_dict(model.state_dict())

        res = 100
        
        if targets is None:
            targets = []
            for i in range(-1 * res, res):
                targets.append(target_model(i / res).data[0])
        
        indexes, outputs = [], []
        for i in range(-1 * res, res):
            indexes.append(i / res)
            outputs.append(model(i / res).data[0])

        fig = plt.figure()
        ax1 = fig.add_subplot(111)

        ax1.plot(indexes, outputs, color='b')
        ax1.plot(indexes, targets, color='g')

        buf = io.BytesIO()
        plt.savefig(buf)

        plt.close()

        return buf.getvalue()
