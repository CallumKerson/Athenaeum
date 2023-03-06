/**
 SearchBar.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct SearchBar: View {
    @Binding var text: String

    @State private var isEditing = false

    var body: some View {
        HStack {
            TextField("Search ...", text: $text)
                .padding(7)
                .padding(.horizontal, 25)
                .background(Color.tertiaryBackground)
                .cornerRadius(8)
                .overlay(
                    HStack {
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.gray)
                            .frame(minWidth: 0, maxWidth: .infinity, alignment: .leading)
                            .padding(.leading, 8)

                        #if os(iOS)
                            if isEditing {
                                Button(action: {
                                    self.text = ""

                                }) {
                                    Image(systemName: "multiply.circle.fill")
                                        .foregroundColor(.gray)
                                        .padding(.trailing, 8)
                                }
                            }
                        #endif
                    }
                )
                .padding(.horizontal, 10)
                .onTapGesture {
                    logger.info("\"Tapped\" on search bar")
                    self.isEditing = true
                }
                .onChange(of: text) { _ in
                    logger.info("\"Tapped\" on search bar")
                    self.isEditing = true
                }

            if isEditing {
                Button(action: {
                    self.isEditing = false
                    self.text = ""

                    #if os(iOS)
                        // Dismiss the keyboard
                        UIApplication.shared.sendAction(
                            #selector(UIResponder.resignFirstResponder),
                            to: nil,
                            from: nil,
                            for: nil
                        )
                    #endif
                }) {
                    Text("Cancel")
                }
                .padding(.trailing, 10)
                .transition(.scale)
                .animation(.default)
            }
        }
    }
}

struct SearchBar_Previews: PreviewProvider {
    static var previews: some View {
        SearchBar(text: .constant(""))
    }
}
