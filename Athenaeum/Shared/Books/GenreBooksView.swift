/**
 GenreBooksView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct GenreBooksView: View {
    @EnvironmentObject private var model: AthenaeumModel
    var genre: Genre

    var body: some View {
        BooksListView(books: model.books.filter { $0.metadata.genre == genre })
            .navigationTitle(genre.rawValue)
    }
}

struct GenreBooksView_Previews: PreviewProvider {
    static var previews: some View {
        GenreBooksView(genre: .sciFi)
            .environmentObject(AthenaeumModel())
    }
}
