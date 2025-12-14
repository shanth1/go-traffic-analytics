import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<div></div>}></Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
