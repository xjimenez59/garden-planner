import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

import 'api_service.dart';

// A screen that allows users to take a picture using a given camera.
class TakePictureScreen extends StatefulWidget {
  const TakePictureScreen({super.key, required this.cameras});
  final List<CameraDescription>? cameras;

  @override
  TakePictureScreenState createState() => TakePictureScreenState();
}

class TakePictureScreenState extends State<TakePictureScreen> {
  late CameraController _controller;
  Future<void>? _initializeControllerFuture;
  late bool uploading = false;

  @override
  void initState() {
    super.initState();
    initCamera(widget.cameras![0]);
  }

  @override
  void dispose() {
    // Dispose of the controller when the widget is disposed.
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Take a picture')),
      body: _controller.value.isInitialized
          ? CameraPreview(_controller)
          : const Center(child: CircularProgressIndicator()),
      floatingActionButton: uploading
          ? const Center(child: CircularProgressIndicator())
          : FloatingActionButton(
              // Provide an onPressed callback.
              onPressed: () async {
                // Take the Picture in a try / catch block. If anything goes wrong,
                // catch the error.
                try {
                  // Ensure that the camera is initialized.
                  await _initializeControllerFuture;

                  // Attempt to take a picture and get the file `image`
                  // where it was saved.
                  final image = await _controller.takePicture();

                  setState(() {
                    uploading = true;
                  });

                  final imageUrl = await uploadImage(image);

                  if (!mounted) return;
                  Navigator.pop(context, imageUrl);
                } catch (e) {
                  // If an error occurs, log the error to the console.
                  print(e);
                }
              },
              child: const Icon(Icons.camera_alt),
            ),
    );
  }

  Future initCamera(CameraDescription cameraDescription) async {
    // create a CameraController
    _controller = CameraController(cameraDescription, ResolutionPreset.medium);
    // Next, initialize the controller. This returns a Future.
    try {
      await _controller.initialize().then((_) {
        if (!mounted) return;
        setState(() {});
      });
    } on CameraException catch (e) {
      debugPrint("camera error $e");
    }
  }
}

Future<String> uploadImage(XFile image) async {
  final imageBytes = await image.readAsBytes();
  String result = await ApiService().postPicture(imageBytes);
  return result;
}
