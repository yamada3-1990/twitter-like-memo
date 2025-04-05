import React, { useState, useRef } from "react";
import { searchMemosByKeyword } from '../api';

export const SearchMemoByKeyword = ({ onSuccess }) => {
    const initialState = {
        keyword: '',
    };

    const [values, setValues] = useState(initialState);
    const [searchResults, setSearchResults] = useState([]);
    const [hasSearched, setHasSearched] = useState(false);
    const keywordRef = useRef(null);

    const onValueChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        });
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        
        // キーワードが空の場合は検索を実行しない
        if (!values.keyword.trim()) {
            alert('キーワードを入力してください');
            return;
        }
        
        setHasSearched(true);

        try {
            const results = await searchMemosByKeyword(values);
            setSearchResults(results);
            if (onSuccess) {
                onSuccess();
            }
        } catch (error) {
            console.error('GET error:', error);
            alert(error.message);
        }
    };

    const handleClear = () => {
        setValues(initialState);
        setSearchResults([]);
        setHasSearched(false);
        if (keywordRef.current) {
            keywordRef.current.value = '';
        }
    };

    return (
        <div>
            <div className="search-input-container">
                <input
                    ref={keywordRef}
                    className='input-tag'
                    type='text'
                    name="keyword"
                    placeholder='キーワード検索'
                    onChange={onValueChange}
                    value={values.keyword}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            handleSubmit(e);
                        }
                    }}
                />
                <button 
                    type="button" 
                    onClick={handleSubmit}
                    className="search-button"
                >
                    検索
                </button>
                {hasSearched && (
                    <button 
                        type="button" 
                        onClick={handleClear}
                        className="clear-button"
                    >
                        X 結果をクリア
                    </button>
                )}
            </div>
            
            <div className="search-results">
                {searchResults.length > 0 ? (
                    searchResults.map((result, index) => (
                        <div key={index} className="timeline">
                            <div className='memo-card'>
                                <div className='memo-content'>
                                    <h3>{result.title}</h3>
                                    <p>{result.body}</p>
                                    {result.tags && (
                                        <div className="memo-tags">
                                            {result.tags.split(',').map((tag, tagIndex) => (
                                                <span key={tagIndex} className="tag">#{tag.trim()}</span>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))
                ) : (
                    hasSearched && <p>検索結果がありません</p>
                )}
            </div>
        </div>
    );
};