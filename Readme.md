# An image processing tool that performs basic operations
## Operations: 
- -grayscale (recursively searches the passed directory and creates grayscale copies of the images)
- -smoothen (creates smoothened copies of the images)
- -sharpen (creates smoothened copies of the images)
- -find=<specific-object> (creates copies of the images that match the <specific-object>)


## TO RUN
- git clone the repo
- docker build -t imagename
- go to the directory you'd like to have your images manipulated and run ```docker run  --rm -it --mount type=bind,source="$(pwd)",target=/app/images imagename <command>```

## Model used for classifying the images is the coco model with pretrained image data, all of which can be found under tensorflowAPI/model/frozen_inference_graph.pb
- the names of the graph operations present in the model are the following:
-  input operation: "image_tensor"
-  output operation: "detection_boxes"
-  output operation: "detection_scores"
-  output operation: "detection_classes"
-  output operation: "num_detections"