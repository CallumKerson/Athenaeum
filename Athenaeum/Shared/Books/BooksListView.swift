/**
 BooksListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct BooksListView: View {
    var books: [Book]

    @State private var selection: Book?
    @EnvironmentObject private var model: AthenaeumModel

    var content: some View {
        List(selection: $selection) {
            ForEach(books) { book in
                NavigationLink(
                    destination: BookTextView(metadata: book.metadata)
                        .environmentObject(model),
                    tag: book,
                    selection: $selection
                ) {
                    BookRowView(book: book)
                }
                .tag(book)
                .onReceive(model.$selectedBookId) { newValue in
                    guard let bookId = newValue, let book = model.books.first(where: { $0.id == bookId }) else { return }
                    selection = book
                }
            }
        }
    }

    @ViewBuilder var body: some View {
        #if os(iOS)
            content
        #else
            content
                .frame(minWidth: 270, idealWidth: 300, maxWidth: 400, maxHeight: .infinity)
                .toolbar { Spacer() }
        #endif
    }
}

struct BooksListView_Previews: PreviewProvider {
    static var previews: some View {
        ForEach([ColorScheme.light, .dark], id: \.self) { scheme in
            NavigationView {
                BooksListView(books: [])
                    .navigationTitle("Books")
                    .environmentObject(AthenaeumModel())
            }
            .preferredColorScheme(scheme)
        }
    }
}
