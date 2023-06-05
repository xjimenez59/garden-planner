import 'dart:convert';
import 'dart:developer';

import 'package:flutter/foundation.dart';
import 'package:flutter_guid/flutter_guid.dart';
import 'package:http/http.dart' as http;
import 'constants.dart';
import 'action_log.dart';
import 'legumes_model.dart';

class ApiService {
  Future<List<ActionLog>?> getLogs() async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.logsEndpoint);
      var response = await http.get(url);
      if (response.statusCode == 200) {
        List<ActionLog> model = actionLogFromJson(response.body);
        return model;
      }
    } catch (e) {
      log(e.toString());
    }
    return [];
  }

  Future<List<Legume>?> getLegumes() async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.legumesEndPoint);
      var response = await http.get(url);
      if (response.statusCode == 200) {
        List<Legume> model = LegumeFromJson(response.body);
        return model;
      }
    } catch (e) {
      log(e.toString());
    }
    return [];
  }

  Future<List<String>?> getTags() async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.tagsEndPoint);
      var response = await http.get(url);
      if (response.statusCode == 200) {
        List<String> result = List<String>.from(jsonDecode(response.body));
        return result;
      }
    } catch (e) {
      log(e.toString());
    }
    return [];
  }

  Future<List<String>?> getLieux() async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.lieuxEndPoint);
      var response = await http.get(url);
      if (response.statusCode == 200) {
        List<String> result = List<String>.from(jsonDecode(response.body));
        return result;
      }
    } catch (e) {
      log(e.toString());
    }
    return [];
  }

  Future<int?> postLogs(List<ActionLog> logs) async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.logsEndpoint);
      var response = await http.post(url, body: jsonEncode(logs));
      if (response.statusCode == 201) {
        dynamic result = jsonDecode(response.body);
        return result["updated"];
      }
    } catch (e) {
      log(e.toString());
    }
    return 0;
  }

  Future<String?> postLog(ActionLog a) async {
    try {
      var url = Uri.parse(ApiConstants.baseUrl + ApiConstants.logEndpoint);
      var response = await http.post(url, body: jsonEncode(a));
      if (response.statusCode == 201) {
        dynamic result = jsonDecode(response.body);
        return result["_id"];
      }
    } catch (e) {
      log(e.toString());
    }
    return null;
  }

  Future<bool> deleteLog(String id) async {
    try {
      var url =
          Uri.parse("${ApiConstants.baseUrl}${ApiConstants.logEndpoint}/$id");
      var response = await http.delete(url);
      if (response.statusCode == 200) {
        dynamic result = jsonDecode(response.body);
        return true;
      }
    } catch (e) {
      log(e.toString());
    }
    return false;
  }

  Future<String> postPicture(Uint8List imageBytes) async {
    try {
      final url = Uri.parse(ApiConstants.baseUrl + ApiConstants.photoEndPoint);
      var request = http.MultipartRequest('POST', url);
      request.files.add(http.MultipartFile.fromBytes('file', imageBytes,
          filename: Guid.newGuid.toString()));

      var streamedResponse = await request.send();
      var response = await http.Response.fromStream(streamedResponse);
      if (response.statusCode == 200) {
        dynamic result = jsonDecode(response.body);
        final String fileUrl =
            "https://storage.googleapis.com${result['pathname']}";
        return fileUrl;
      }
    } catch (e) {
      log(e.toString());
    }
    return "";
  }
}
