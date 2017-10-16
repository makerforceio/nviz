import torch

import torch.nn as nn

from torch.nn import Parameter
from torch.autograd import Variable

class Model(nn.Module):
    def __init__(self, polySize):
        super(Model, self).__init__()

        self.wx = Parameter(torch.randn(polySize), requires_grad=True)
        self.polySize = polySize
        
    def forward(self, x):
        y = 0
        xe = 1
        for i in range(0, self.polySize):
            y += self.wx[i] * xe
            xe = xe * x

        return y

