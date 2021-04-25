/**
 Sidebar.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

//
//  Sidebar.swift
//  Athenaeum (macOS)
//
//  Created by Callum Kerson on 21/04/2021.
//
import SwiftUI

struct Sidebar: View {
    @Binding var selectedBookId: String?

    var body: some View {
        List(Genre.allCases) { genre in
            NavigationLink(destination: BooksListView(
                viewModel: BooksListViewModel(genre: genre),
                selectedBookId: $selectedBookId
            )) {
                Label(genre.name, systemImage: genre.icon)
            }
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 150, idealWidth: 150, maxWidth: 200, maxHeight: .infinity)
        .padding(.top, 16)
    }
}
