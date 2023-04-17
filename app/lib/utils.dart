String dateFormat(DateTime date) {
  var weekdays = ["oups", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];
  var monthNames = [
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
  String strDate =
      "${weekdays[date.weekday]} ${date.day} ${monthNames[date.month]} ${date.year}";
  return strDate;
}

int weekNum(DateTime date) {
  int weeknum =
      (date.difference(DateTime.utc(date.year, 1, 1)).inDays / 7).ceil() + 1;
  return weeknum;
}
