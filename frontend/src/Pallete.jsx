import { useCallback, useEffect, useState } from 'react';
import DraggableItem from './DraggableItem.jsx';

const Pallete = ({ measurementList, fieldsList, bucketList}) => {
    return (
      <div className="p-4 w-64 h-screen overflow-y-auto text-slate-100">
        {/*<h3 className="text-lg font-bold mb-4 text-slate-100">Buckets</h3>
        <div className="space-y-2">
        {/*bucketList.map((bucket, index) => (
          <DraggableItem key={index} item={bucket.name} type="BUCKET" />
        ))}
        </div>*/}
        <h3 className="text-lg font-bold mb-4 mt-6 text-slate-100">Measurements</h3>
        <div className="space-y-2 text-slate-100">
        {measurementList.map((measurement, index) => (
          <DraggableItem key={index} item={measurement} type="MEASUREMENT" />
        ))}
        </div>
        <h3 className="text-lg font-bold mb-4 mt-6 text-slate-100">Fields</h3>
        <div className="space-y-2 text-slate-100">
        {fieldsList.map((field, index) => (
          <DraggableItem key={index} item={field} type="FIELD" />
        ))}
        </div>
      </div>
    );
  };
  
  export default Pallete;