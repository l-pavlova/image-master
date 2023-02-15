FROM golang:1.20

# 1. Install the TensorFlow C Library (v2.11.0).
RUN curl -L https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-$(uname -m)-2.11.0.tar.gz \
    | tar xz --directory /usr/local \
    && ldconfig

# 2. Install the Protocol Buffers Library and Compiler.
RUN apt-get update && apt-get -y install --no-install-recommends \
    libprotobuf-dev \
    protobuf-compiler
RUN apt-get update && apt-get -y install unzip
# 3. Install and Setup the TensorFlow Go API.
RUN git clone --branch=v2.11.0 https://github.com/tensorflow/tensorflow.git /go/src/github.com/tensorflow/tensorflow \
    && cd /go/src/github.com/tensorflow/tensorflow \
    && go mod init github.com/tensorflow/tensorflow \
    && sed -i '72 i \    ${TF_DIR}\/tensorflow\/tsl\/protobuf\/*.proto \\' tensorflow/go/genop/generate.sh \
    && (cd tensorflow/go/op && go generate) \
    && go mod edit -require github.com/google/tsl@v0.0.0+incompatible \
    && go mod edit -replace github.com/google/tsl=/go/src/github.com/google/tsl \
    && (cd /go/src/github.com/google/tsl && go mod init github.com/google/tsl && go mod tidy) \
    && go mod tidy \
    && go test ./...

# #End: install tensorflow
# RUN mkdir -p /model && \
#   curl -o /model/inception5h.zip -s "http://download.tensorflow.org/models/inception5h.zip" && \
#   unzip /model/inception5h.zip -d /model


WORKDIR /app
COPY . .

RUN ls -R
RUN pwd


# Install the app
RUN go build -o imagemaster
RUN go test

ENTRYPOINT [ "/app/imagemaster" ]