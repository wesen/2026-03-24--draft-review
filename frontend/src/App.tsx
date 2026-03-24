import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthorApp } from "./app/AuthorApp";
import { ReaderApp } from "./app/ReaderApp";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/r/:token" element={<ReaderApp />} />
        <Route path="/*" element={<AuthorApp />} />
      </Routes>
    </BrowserRouter>
  );
}
