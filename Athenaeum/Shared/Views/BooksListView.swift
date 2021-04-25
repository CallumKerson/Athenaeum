/**
 BooksListView.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import SwiftUI

struct BooksListView: View {
    @ObservedObject var viewModel: BooksListViewModel
    @Binding var selectedBookId: String?

    var body: some View {
        #if os(macOS)
            return view
                .frame(minWidth: 400, minHeight: 600)
                .toolbar {
                    ToolbarItem(placement: .automatic) {
                        Button(action: { viewModel.reload() }) {
                            Image(systemName: "arrow.clockwise")
                        }
                        .keyboardShortcut("R", modifiers: .command)
                    }
                    ToolbarItem(placement: .automatic) {
                        Button(action: {
                            NSApp.keyWindow?.firstResponder?.tryToPerform(
                                #selector(NSSplitViewController.toggleSidebar(_:)),
                                with: nil
                            )
                        }) {
                            Image(systemName: "sidebar.left")
                        }
                        .keyboardShortcut("S", modifiers: .command)
                    }
                }
        #else
            return view
                .toolbar {
                    ToolbarItem(placement: .navigationBarTrailing) {
                        Button(action: { viewModel.reload() }) {
                            Image(systemName: "arrow.clockwise")
                        }
                    }
                }
        #endif
    }

    @ViewBuilder
    private var view: some View {
        if viewModel.loading {
            ProgressView()
                .onAppear(perform: { viewModel.reload() })
                .frame(maxWidth: .infinity, maxHeight: .infinity)
        } else if let error = viewModel.error {
            Label(error.description, systemImage: "exclamationmark.triangle")
        } else {
            List(viewModel.items) { item in
                NavigationLink(
                    destination: BookDetailView(viewModel: BookDetailViewModel(id: item.id)),
                    tag: item.id,
                    selection: $selectedBookId
                ) {
                    BookCellView(viewModel: BookCellViewModel(id: item.id))
                }
            }
        }
    }
}
