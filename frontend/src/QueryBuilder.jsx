import { useCallback, useEffect, useState } from 'react';
import DropSelector from './DropSelector';
import DropSelectorContainer from './DropSelectorContainer';
import { DateTimePicker } from '@atlaskit/datetime-picker';

function useBuckets() {
    const [bucketList, setBucketList] = useState([]);
    const [selectedBucket, setSelectedBucket] = useState('');
  
    useEffect(() => {
      async function fetchBuckets() {
        try {
          const response = await fetch('/get-buckets');
          if (!response.ok) {
            throw new Error('Failed to fetch buckets');
          }
          const buckets = await response.json();
          setBucketList(buckets);
          setSelectedBucket(buckets[0]?.name || '');
        } catch (error) {
          console.error('Error fetching buckets:', error);
        }
      }
      fetchBuckets();
      createMap();
    }, []);
  
    return [bucketList, selectedBucket, setSelectedBucket];
  }
  
 
  function useMeasurements(selectedBucket) {
    const [measurementList, setMeasurementList] = useState([]);
    const [selectedMeasurement, setSelectedMeasurement] = useState('');
  
    useEffect(() => {
      if (!selectedBucket) return;
  
      async function fetchMeasurements() {
        try {
          const response = await fetch(`/get-measurements?bucket=${selectedBucket}`);
          if (!response.ok) {
            throw new Error('Failed to fetch measurements');
          }
          const data = await response.json();
          setMeasurementList(data.measurements);
          setSelectedMeasurement(data.measurements[0] || '');
        } catch (error) {
          console.error('Error fetching measurements:', error);
        }
      }
  
      fetchMeasurements();
    }, [selectedBucket]);
  
    return [measurementList, selectedMeasurement, setSelectedMeasurement];
  }
  

  function useFields(selectedBucket, selectedMeasurement) {
    const [fieldList, setFieldList] = useState([]);
    const [selectedField, setSelectedField] = useState('');
  
    useEffect(() => {
      if (!selectedBucket || !selectedMeasurement) return;
  
      async function fetchFields() {
        try {
          const response = await fetch(
            `/get-fields?bucket=${selectedBucket}&measurement=${selectedMeasurement}`
          );
          if (!response.ok) {
            throw new Error('Failed to fetch fields');
          }
          const data = await response.json();
          setFieldList(data.fields);
          setSelectedField(data.fields[0] || '');
        } catch (error) {
          console.error('Error fetching fields:', error);
        }
      }
  
      fetchFields();
    }, [selectedBucket, selectedMeasurement]);
  
    return [fieldList, selectedField, setSelectedField];
  }

const sendQueryJSON = async (queryJSON) => {
    try {
        console.log(JSON.stringify(queryJSON, null, 2))
        const response = await fetch('/run-query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json', // specify the content type
            },
            body: JSON.stringify(queryJSON, null, 2), // the body must be serialized into JSON format
        });
        console.log('Response:', await response.json());
    } catch (error) {
        console.error('Error:', error);
    }
};

async function createMap(bucket, setFieldMap) {
    try {
        // gets the list of measurements
        const response = await fetch(`/get-measurements?bucket=${bucket}`);
        const measurements = (await response.json())["measurements"];

        // gets the list of fields asynchronously
        const result = {};
        await Promise.all(
            measurements.map(async (measurement) => {
                const response = await fetch(`/get-fields?bucket=${bucket}&measurement=${measurement}`);
                const fields = (await response.json())["fields"];
                result[measurement] = fields;
            })
        );
        // returns the results
        setFieldMap(result);
    } catch (error) {
        console.error("Error fetching measurements and/or fields:", error);
    }
}

export default function QueryBuilder({setExternalTimeRange, setExternalBucketList, setExternalMeasurementList, setExternalFieldList}) {

    const [bucketList, selectedBucket, setSelectedBucket] = useBuckets();
    const [measurementList, selectedMeasurement, setSelectedMeasurement] = useMeasurements(selectedBucket);
    const [fieldList, selectedField, setSelectedField] = useFields(selectedBucket, selectedMeasurement);

    const [timeRange, setTimeRange] = useState({ start: '', end: '' });
    const [filterList, setFilters] = useState([]);
    const [formData, setFormData] = useState({});
    const [fieldMap, setFieldMap] = useState({});

    
    
    const operatorList = ["==", "<=", ">=", "<", ">"];

    useEffect(() => {
        setExternalMeasurementList(measurementList);
    }, [measurementList, setExternalMeasurementList]);

    useEffect(() => {
        setExternalFieldList(fieldList);
    }, [fieldList, setExternalFieldList]);


    useEffect(() => {
        setExternalBucketList(bucketList);
    }, [bucketList, setExternalBucketList]);
    
    useEffect(() => {
        async function createMap() {
          if (!selectedBucket) return;
          try {
            const response = await fetch(`/get-measurements?bucket=${selectedBucket}`);
            const measurements = (await response.json()).measurements;
    
            const result = {};
            await Promise.all(
              measurements.map(async (measurement) => {
                const response = await fetch(
                  `/get-fields?bucket=${selectedBucket}&measurement=${measurement}`
                );
                const fields = (await response.json()).fields;
                result[measurement] = fields;
              })
            );
            setFieldMap(result);
          } catch (error) {
            console.error('Error fetching measurements and/or fields:', error);
          }
        }
    
        createMap();
      }, [selectedBucket]);

      useEffect(() => {
        changeFilterBucket();
      }, [fieldMap]);

      useEffect(() => {
        setFormData((prevData) => ({ ...prevData, bucket: selectedBucket }));
      }, [selectedBucket, formData]);
    
      useEffect(() => {
        setFormData((prevData) => ({ ...prevData, measurements: selectedMeasurement }));
      }, [selectedMeasurement, formData]);
    
      useEffect(() => {
        setFormData((prevData) => ({ ...prevData, fields: selectedField }));
      }, [selectedField, formData]);
    
      const changeBucket = useCallback(
        (event) => {
          const newBucket = event.target.value;
          setSelectedBucket(newBucket);
        },
        [setSelectedBucket]
      );
    
      const changeMeasurement = useCallback(
        (event) => {
          const newMeasurement = event.target.value;
          setSelectedMeasurement(newMeasurement);
        },
        [setSelectedMeasurement]
      );
    
      const changeField = useCallback(
        (event) => {
          const newField = event.target.value;
          setSelectedField(newField);
        },
        [setSelectedField]
      );
    
      const changeTimeRange = useCallback(({ name, value }) => {
        setTimeRange((prevRange) => ({ ...prevRange, [name]: value }));
        setFormData((prevData) => ({
          ...prevData,
          timeRange: { ...prevData.timeRange, [name]: value },
        }));
      }, []);
    

      // START FILTERS
      const addFilter = () => {
          const newMeasurement = measurementList.length > 0 ? measurementList[0] : "";
          const newField = measurementList.length > 0 ? fieldMap[newMeasurement][0] : "";

          const newFilter = { measurement: newMeasurement, field: newField, operator: '==', value: '' };
          setFilters(prevFilters => [...prevFilters, newFilter]);
          setFormData(prevData => ({
              ...prevData,
              filters: [...(prevData.filters || []), newFilter]
          }));
      }

      const changeFilterBucket = useCallback(() => {
          const newMeasurement = measurementList.length > 0 ? measurementList[0] : "";
          const newField = measurementList.length > 0 ? fieldMap[newMeasurement] !== undefined ? fieldMap[newMeasurement][0] : "" : "";

          setFilters(prevFilters => prevFilters !== undefined ?
              prevFilters.map((filter) => 
                  ({ ...filter, measurement: newMeasurement, field: newField})
              ) : {}
          );
          setFormData(prevData => prevData.filters !== undefined ? (  {
              ...prevData,
              filters: prevData.filters.map((filter) =>
                  ({ ...filter, measurement: newMeasurement, field: newField })
              )
          }) : {}
          );
      }, [measurementList, fieldMap, setFilters, setFormData]);

      const removeFilter = (index) => {
          setFilters(prevFilters => prevFilters.filter((_, i) => i !== index));
          setFormData(prevData => ({
              ...prevData,
              filters: prevData.filters.filter((_, i) => i !== index)
          }));
      }

      const changeFilter = useCallback((index, key, value) => {
          setFilters(prevFilters =>
              prevFilters.map((filter, i) =>
                  i === index ? { ...filter, [key]: value } : filter
              )
          );
          setFormData(prevData => ({
              ...prevData,
              filters: prevData.filters.map((filter, i) =>
                  i === index ? { ...filter, [key]: value } : filter
              )
          }));
      }, []);

    // END FILTERS

    const handleSubmit = useCallback(() => {
        console.log(timeRange)
        setExternalTimeRange({from: Date.parse(timeRange.start), to: Date.parse(timeRange.stop), state: Math.random()})
        console.log({from: Date.parse(timeRange.start), to: Date.parse(timeRange.stop)})
        sendQueryJSON(formData)

        const jsonOutput = JSON.stringify(formData, null, 2);
        console.log(jsonOutput);
        // You can also send this JSON to a server or use it as needed
    }, [formData, timeRange, setExternalTimeRange]);

    return (
        <div className="QueryBuilder">
            <div style={{ width: '100%', height: '100%' }}>
                <div className='bucket-selector-container'>
                    <label>Bucket:</label>
                    <select value={selectedBucket} onChange={changeBucket}>
                        {bucketList.map((bucket) => (
                            <option key={bucket.name} value={bucket.name}>{bucket.name}</option>
                        ))}
                    </select>
                </div>
                <DropSelectorContainer Type="MEASUREMENT" setExternalSelectList={setSelectedMeasurement} externalSelectList={selectedMeasurement} onChange={changeMeasurement} title='Measurments' dropText={'drag measurements here'} handleSubmit={handleSubmit} />
                <DropSelectorContainer Type="FIELD" setExternalSelectList={setSelectedField} externalSelectList={selectedField} onChange={changeField} title='Fields' dropText={'drag fields here'} handleSubmit={handleSubmit} />
                <div className='range-selector-container'>
                    <label>Range:</label>
                    <label htmlFor="start">Start:</label>
                    <DateTimePicker
                        name="start"
                        value={timeRange.start}
                        onChange={(value) => changeTimeRange({ name: 'start', value })}
                    />
                    <label htmlFor="stop">Stop:</label>
                    <DateTimePicker
                        name="stop"
                        value={timeRange.stop}
                        onChange={(value) => changeTimeRange({ name: 'stop', value })}
                    />
                </div>
                <div>
                    <div>
                        <button onClick={() => addFilter()} className="w-full">
                            Add Filter
                        </button>
                    </div>
                    <div>
                        {filterList.map((filter, index) => (
                            <div key={index}>
                                <select
                                    value={filter.measurement}
                                    onChange={(e) => changeFilter(index, 'measurement', e.target.value)}>
                                    {measurementList.map((measurement) => (
                                        <option key={measurement} value={measurement}>{measurement}</option>
                                    ))}
                                </select>

                                <select
                                    value={filter.field}
                                    onChange={(e) => changeFilter(index, 'field', e.target.value)}>
                                    {fieldMap[filter.measurement] ? fieldMap[filter.measurement].map((field) => (
                                        <option key={field} value={field}>{field}</option>
                                    )) : <></>}
                                </select>
                                <select value={filter.operator}
                                    onChange={(e) => changeFilter(index, 'operator', e.target.value)}
                                    default="==">
                                    {operatorList.map((operator) => (
                                        <option key={operator} value={operator}>{operator}</option>
                                    ))}
                                </select>
                                <input
                                    type="text"
                                    value={filter.value}
                                    onChange={(e) => changeFilter(index, 'value', e.target.value)}
                                    placeholder="Threshold"
                                />

                                <button onClick={() => removeFilter(index)}>Remove</button>
                            </div>
                        ))}
                    </div>
                </div>
                <button onClick={handleSubmit}>Submit</button>
            </div>
        </div>
    );
}