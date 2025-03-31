import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  return (
    <>
      <div className='add-memo'>
        <input className='input-title' type='text' placeholder='タイトル' />
        <input className='input-body' type='text' placeholder='本文' />
        <input className='input-tag' type='text' placeholder='タグ' />
        <button className='post-button'>メモを追加</button>
      </div>

      <div className='search'>
        <form action="検索">
          <input className='search-keyword' type='text' placeholder='キーワード検索' />
          <input className='search-tag' type='text' placeholder='タグ検索' />
          <button className='search-button'>検索</button>
        </form>
      </div>

      <div className='timeline'>
        <div className='memo-card'>
          <div className='memo-title'>テストタイトル1</div>
          <div className='memo-body'>テスト本文1</div>
          {/* <div>1990-02-20 wed</div> */}
        </div>

        <hr />

        <div className='memo-card'>
          <div className='memo-title'>テストタイトル2</div>
          <div className='memo-body'>テスト本文2</div>
          {/* <div>1990-02-20 wed</div> */}
        </div>

        <hr />

        <div className='memo-card'>
          <div className='memo-title'>テストタイトル2</div>
          <div className='memo-body'>テスト本文2</div>
          {/* <div>1990-02-20 wed</div> */}
        </div>
      </div>
    </>
  )
}

export default App
