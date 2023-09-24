const weekdays = ["oups", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];
const monthNames = [
  "ouch",
  "Janvier",
  "Février",
  "Mars",
  "Avril",
  "Mai",
  "Juin",
  "Juillet",
  "Août",
  "Septembre",
  "Octobre",
  "Novembre",
  "Décembre"
];

String dateFormat(DateTime date) {
  String strDate =
      "${weekdays[date.weekday]} ${date.day} ${monthNames[date.month]} ${date.year}";
  return strDate;
}

String weightFormat(int? g) {
  String res = "";
  if (g == null) return "";
  if (g >= 1000) {
    double kg = g / 1000;
    res = "${kg.toStringAsFixed(1)} kg";
  } else if (g > 0) {
    res = "$g g";
  }
  return res;
}

int weekNum(DateTime date) {
  int weeknum =
      (date.difference(DateTime.utc(date.year, 1, 1)).inDays / 7).ceil() + 1;
  return weeknum;
}

extension DateTools on DateTime {
  bool sameDayAs(DateTime other) {
    return (year == other.year) && (month == other.month) && (day == other.day);
  }
}

extension DiacriticsAwareString on String {
  static const diacritics =
      'ÀÁÂÃÄÅàáâãäåÒÓÔÕÕÖØòóôõöøÈÉÊËĚèéêëěðČÇçčÐĎďÌÍÎÏìíîïĽľÙÚÛÜŮùúûüůŇÑñňŘřŠšŤťŸÝÿýŽž';
  static const nonDiacritics =
      'AAAAAAaaaaaaOOOOOOOooooooEEEEEeeeeeeCCccDDdIIIIiiiiLlUUUUUuuuuuNNnnRrSsTtYYyyZz';

  String get withoutDiacriticalMarks => splitMapJoin('',
      onNonMatch: (char) => char.isNotEmpty && diacritics.contains(char)
          ? nonDiacritics[diacritics.indexOf(char)]
          : char);
}
