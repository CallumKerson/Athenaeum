/**
 LibraryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import RealmSwift
import SwiftUI

struct LibraryView: View {
    @State var books: [Audiobook]

    var body: some View {
        NavigationView {
            List {
                Section(header: HeaderView()) {
                    ForEach(books) { book in
                        NavigationLink(destination: AudiobookView(book: book)) {
                            AudiobookCellView(book: book)
                        }
                    }
                }
            }
            .frame(minWidth: 425, maxWidth: 425)
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 850, maxWidth: 850, minHeight: 400, maxHeight: .infinity)
//        .onAppear {
//            self.loadBooks()
//        }
    }

//    func loadBooks() {
//        books = Library.global.audiobooks
//    }
}

struct HeaderView: View {
    var body: some View {
        VStack {
            HStack(spacing: 20) {
                Text("Library")
                    .layoutPriority(1)
                    .font(.largeTitle)
                Spacer()
            }
        }.padding(.bottom)
    }
}

struct LibraryView_Previews: PreviewProvider {
    static var previews: some View {
        LibraryView(books: [
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/TWoK.m4b"),
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/In the Labyrinth of Drakes.m4b"),
            Audiobook(fromFileWithPath: "/Users/ckerson/Music/Smarter Faster Better.m4b"),
        ])
    }
}
