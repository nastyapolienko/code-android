import 'package:flutter/material.dart';
import 'package:fetchlist/main.dart';
import 'dart:async';
import 'package:http/http.dart' as http;
import 'dart:convert';

Future<Book> fetchBook() async {
  //var a = index.toString();
  final response = await http.get('http://192.168.0.5:8080/books/1');

  if (response.statusCode == 200) {
    return Book.fromJson(json.decode(response.body));
  } else {
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
          title: Text(book.bookname!=null?book.bookname:'The details'),
        ),
        body: Center(

          child: FutureBuilder<Book>(

            future: fetchBook(),
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                return Text(snapshot.data.bookname + "\n" + snapshot.data.year);
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
  final String year;

  Book({this.id, this.bookname, this.year});

  factory Book.fromJson(Map<String, dynamic> json) {
    return Book(
      id: json['id'],
      bookname: json['bookname'],
      year: json['year'],
    );
  }
}
