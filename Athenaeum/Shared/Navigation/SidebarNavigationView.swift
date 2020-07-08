/**
 SidebarNavigationView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct SidebarNavigationView: View {
    enum NavigationItem {
        case all
        case genre
    }

    @EnvironmentObject private var model: AthenaeumModel
    @State private var selections: Set<NavigationItem> = [.all]
    @State private var presentingRewards = false

    var sidebar: some View {
        List(selection: $selections) {
            NavigationLink(destination: AllBooksView()) {
                Label("All Books", systemImage: "books.vertical")
            }
            .accessibility(label: Text("All Books"))
            .tag(NavigationItem.all)
            NavigationLink(destination: GenreBooksView(genre: .sciFi)) {
                Label("Sci-Fi", systemImage: "bolt")
            }
            .accessibility(label: Text("Sci-Fi"))
            .tag(NavigationItem.genre)
        }
        .listStyle(SidebarListStyle())
    }

    var body: some View {
        NavigationView {
            #if os(macOS)
                sidebar.frame(minWidth: 100, idealWidth: 150, maxWidth: 200, maxHeight: .infinity)
            #else
                sidebar
            #endif

            Text("Content List")
                .frame(maxWidth: .infinity, maxHeight: .infinity)

            #if os(macOS)
                Text("Select a Book")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
                    .toolbar { Spacer() }
            #else
                Text("Select a Book")
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            #endif
        }
    }
}

struct SidebarNavigationView_Previews: PreviewProvider {
    static var previews: some View {
        SidebarNavigationView()
    }
}
