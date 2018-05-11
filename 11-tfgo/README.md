### Understanding Tensorflow using Go

Installing TensorFlow for Go: https://www.tensorflow.org/versions/master/install/install_go

Install TensorFlow on Mac: https://www.tensorflow.org/versions/master/install/install_mac#determine_which_tensorflow_to_install

```
docker pull tensorflow/tensorflow
```

Install C lib:
```
TF_TYPE="cpu"
TARGET_DIRECTORY='/usr/local'
curl -L \
"https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-1.7.0-rc1.tar.gz" |
sudo tar -C $TARGET_DIRECTORY -xz

export LIBRARY_PATH=$LIBRARY_PATH:$TARGET_DIRECTORY
```