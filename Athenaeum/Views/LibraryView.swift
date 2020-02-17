/**
 LibraryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import RealmSwift
import SwiftUI

struct LibraryView: View {
    @State private var books: Results<Audiobook> = Library.global.audiobooks

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
            .frame(minWidth: 250, maxWidth: 350)
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 400, maxWidth: .infinity, minHeight: 400, maxHeight: .infinity)
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
        HStack(spacing: 20) {
            Text("Library")
                .layoutPriority(1)
                .font(.largeTitle)
            Spacer()
        }
    }
}

struct LibraryView_Previews: PreviewProvider {
    static var previews: some View {
        LibraryView()
    }
}
