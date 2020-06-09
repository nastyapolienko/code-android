import 'package:flutter/material.dart';
import 'package:fetchlist/main.dart';
import 'dart:async';
import 'package:http/http.dart' as http;
import 'dart:convert';

Future<Book> fetchBook() async {
  final response =
  await http.get('http://192.168.0.108:8080/books/1');

  if (response.statusCode == 200) {
    // If the server did return a 200 OK response,
    // then parse the JSON.
    return Book.fromJson(json.decode(response.body));
  } else {
    // If the server did not return a 200 OK response,
    // then throw an exception.
    throw Exception('Failed to load a book');
  }
}



class DetailPage extends StatelessWidget {

  final Book book;
  DetailPage(this.book);


  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(book.bookname),
        ),
        body: Center(
          child: FutureBuilder<Book>(
            future: fetchBook(),
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                return Text(snapshot.data.bookname);
              } else if (snapshot.hasError) {
                return Text("${snapshot.error}");
              }

              // By default, show a loading spinner.
              return CircularProgressIndicator();
            },
          ),
        ),
    );
  }
}

class Book {
  final int id;
  final String bookname;

  Book({this.id, this.bookname});

  factory Book.fromJson(Map<String, dynamic> json) {
    return Book(
      id: json['id'],
      bookname: json['bookname'],
    );
  }
}
