FROM serverwentdown/miniconda

RUN /opt/conda/bin/conda install pytorch torchvision cuda80 matplotlib -c soumith -c conda-forge -y --quiet

COPY . .

ENTRYPOINT ["/opt/conda/bin/python", "wrapper.py"]
