const _apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8081',
);

class ApiConstants {
  static String baseUrl = _apiBaseUrl;

  static String logsEndpoint = '/logs';
  static String logEndpoint = '/log';
  static String legumesEndPoint = '/legumes';
  static String tagsEndPoint = '/tags';
  static String lieuxEndPoint = '/lieux';
  static String photoEndPoint = '/photo';
  static String recoltesEndPoint = '/recoltes';
  static String recolteAnnuelleEndPoint = '/recoltes/lieux';
  static String jardinsEndPoint = '/gardens';
  static String jardinEndPoint = '/garden';
  static String GOOGLE_CLIENT_ID =
      '490039520157-pudl2tcbpsc3sru9ci7caqqtuhakctlf.apps.googleusercontent.com';

  static String defaultUser = 'j2U1HTvFWkOZCBcrmidKcltf0UN2'; // GJ
}
