import 'dart:developer';

import 'package:http/http.dart' as http;
import 'constants.dart';
import 'action_log.dart';

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
  }
}
