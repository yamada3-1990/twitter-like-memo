import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { Adding } from './components/Adding';
import { MemoList } from './components/MemoList';
import { SearchMemoByKeyword } from './components/SearchByKeyword';

function App() {
  const [reload, setReload] = useState(false);

  const handleMemoAdded = () => {
    setReload(prev => !prev);
  };

  const handleMemoDeleted = () => {
    setReload(prev => !prev);
  };

  const handleMemoSearchKeyword = () => {
    setReload(prev => !prev);
  };

  return (
    <>
      <div>
        <Adding onSuccess={handleMemoAdded} />
      </div>
      <div className='search'>
        <form action="検索">
          {/* <input className='search-keyword' type='text' placeholder='キーワード検索' /> */}
          <SearchMemoByKeyword onSuccess={handleMemoSearchKeyword} />
          <input className='search-tag' type='text' placeholder='タグ検索' />
          <button className='search-button'>検索</button>
        </form>
      </div>
      <MemoList reload={reload} onDelete={handleMemoDeleted} />
    </>
  )
}

export default App
