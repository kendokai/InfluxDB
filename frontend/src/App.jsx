import { useEffect, useState } from 'react';
import QueryBuilder from './QueryBuilder';
import Pallete from './Pallete';

import Grafanapanel from './Grafanapanel';
import './index.css';


export default function App() {
  const [measurementList, setMeasurementList] = useState([]);
  const [fieldList, setFieldList] = useState([]);

  const [bucketList, setBucketList] = useState([]);
  const [timeRange, setTimeRange] = useState([]);

  return (
    <>
      <div className="bg-[#252525] fixed top-0 left-0 z-10 w-[270px] h-full p-2.5 flex flex-col items-center overflow-y-auto overflow-x-hidden scrollbar-overlay">
        <Pallete bucketList={bucketList} measurementList={measurementList} fieldsList={fieldList} />
      </div>
      <div className="absolute top-0 right-0 bottom-0 left-[270px] p-4">
        <QueryBuilder setExternalTimeRange={setTimeRange} setExternalBucketList={setBucketList} setExternalMeasurementList={setMeasurementList} setExternalFieldList={setFieldList} />
        <h1>Graph:</h1>
        <Grafanapanel timeRange={timeRange}/>

      </div>
    </>
  );
}
