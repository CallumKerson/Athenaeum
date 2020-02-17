/**
 ContentView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct ContentView: View {
    let books: [Audiobook]

    var body: some View {
        VStack(alignment: .leading) {
            Text("Library")
                .font(.largeTitle)
            AudiobookListView(books: books)
                .shadow(radius: 10)
        }
        .padding()
        .background(Color(NSColor.windowBackgroundColor))
    }
}

struct ContentView_Previews: PreviewProvider {
    static let books = [
        Audiobook.getBookFromFile(path: "/Users/ckerson/Music/TWoK.m4b"),
        Audiobook.getBookFromFile(path: "/Users/ckerson/Music/The Gift.m4b"),
        Audiobook.getBookFromFile(path: "/Users/ckerson/Music/Gothe F--ktoSleep_ep6.m4b"),
    ]
    static var previews: some View {
        ContentView(books: books)
            .environment(\.colorScheme, .light)
//        Group {
//            ContentView(books: books)
//                .environment(\.colorScheme, .light)
//
//            ContentView(books: books)
//                .environment(\.colorScheme, .dark)
//        }
    }
}
