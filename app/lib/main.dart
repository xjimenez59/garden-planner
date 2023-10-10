import 'package:flutter/material.dart';

import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:firebase_core/firebase_core.dart';
import 'firebase_options.dart';
import 'package:firebase_auth/firebase_auth.dart' hide EmailAuthProvider;
import 'package:firebase_ui_auth/firebase_ui_auth.dart';
import 'package:firebase_ui_oauth_google/firebase_ui_oauth_google.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:firebase_ui_localizations/firebase_ui_localizations.dart';

import 'home.dart';

Future<void> main() async {
// Ensure that plugin services are initialized so that `availableCameras()`
  // can be called before `runApp()`
  WidgetsFlutterBinding.ensureInitialized();

  // Obtain a list of the available cameras on the device.
//  final cameras = await availableCameras();

  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );

  FirebaseUIAuth.configureProviders([
    GoogleProvider(
        clientId:
            '490039520157-pudl2tcbpsc3sru9ci7caqqtuhakctlf.apps.googleusercontent.com'),
    EmailAuthProvider()
  ]);

  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  MyApp({super.key});

  /*  final providers = [
    EmailAuthProvider(),
    GoogleAuthProvider(),
  ]; */

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Garden Planner',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      locale: Locale("fr"),
      localizationsDelegates: [
        FirebaseUILocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: [
        Locale('fr'), // French
        Locale('en'), // English
        Locale('es'), // Spanish
      ],
      initialRoute:
          FirebaseAuth.instance.currentUser == null ? '/sign-in' : '/',
      routes: {
        '/sign-in': SignIn,
        '/profile': (context) {
          return ProfileScreen(
            //   providers: providers,
            actions: [
              SignedOutAction((context) {
                Navigator.pushReplacementNamed(context, '/sign-in');
              }),
            ],
          );
        },
        '/': (context) {
          if (FirebaseAuth.instance.currentUser == null) {
            return SignIn(context);
          } else {
            return const MyHomePage(title: 'Garden Planner Home');
          }
        }
      },
    );
  }

  Widget SignIn(context) {
    return SignInScreen(
      //     providers: providers,
      actions: [
        AuthStateChangeAction<SignedIn>((context, state) {
          Navigator.pushReplacementNamed(context, '/');
        }),
      ],
      headerBuilder: (context, constraints, shrinkOffset) {
        return AppBar(
          title: Text("Garden Planner"),
          leading: Icon(Icons.login),
        );
      },
    );
  }
}
