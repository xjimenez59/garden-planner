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

int weekNum(DateTime date) {
  int weeknum =
      (date.difference(DateTime.utc(date.year, 1, 1)).inDays / 7).ceil() + 1;
  return weeknum;
}

extension DateTools on DateTime {
  bool sameDayAs(DateTime other) {
    return (this.year == other.year) &&
        (this.month == other.month) &&
        (this.day == other.day);
  }
}

extension DiacriticsAwareString on String {
  static const diacritics =
      'ÀÁÂÃÄÅàáâãäåÒÓÔÕÕÖØòóôõöøÈÉÊËĚèéêëěðČÇçčÐĎďÌÍÎÏìíîïĽľÙÚÛÜŮùúûüůŇÑñňŘřŠšŤťŸÝÿýŽž';
  static const nonDiacritics =
      'AAAAAAaaaaaaOOOOOOOooooooEEEEEeeeeeeCCccDDdIIIIiiiiLlUUUUUuuuuuNNnnRrSsTtYYyyZz';

  String get withoutDiacriticalMarks => this.splitMapJoin('',
      onNonMatch: (char) => char.isNotEmpty && diacritics.contains(char)
          ? nonDiacritics[diacritics.indexOf(char)]
          : char);
}
