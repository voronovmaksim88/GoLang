<?php

include_once 'segnetics.php';


// out
$Dev1_proc = 0;
$Dev2_proc = 0;
$Dev3_proc = 0;
$Dev4_proc = 0;
$Dev5_proc = 0;


$seg = new SegShm();
$Work = 0; // Флаг режима Работа
$Start_End = false; //флаг окончания первоначальной установки процентов устройств
$Dev1_Isost = 0; //интегральная составляющая регулятора 1-ого устройства
$Dev2_Isost = 0; //интегральная составляющая регулятора 2-ого устройства
$Dev3_Isost = 0; //интегральная составляющая регулятора 3-ого устройства
$Dev4_Isost = 0; //интегральная составляющая регулятора 4-ого устройства
$Dev5_Isost = 0; //интегральная составляющая регулятора 5-ого устройства

$Dev1_Psost = 0; //пропорциональная составляющая регулятора 1-ого устройства
$Dev2_Psost = 0; //пропорциональная составляющая регулятора 2-ого устройства
$Dev3_Psost = 0; //пропорциональная составляющая регулятора 3-ого устройства
$Dev4_Psost = 0; //пропорциональная составляющая регулятора 4-ого устройства
$Dev5_Psost = 0; //пропорциональная составляющая регулятора 5-ого устройства

$Dev1_old = 0; // старая роль устройства для  отслеживания изменения роли на ходу
$Dev2_old = 0; // старая роль устройства для  отслеживания изменения роли на ходу
$Dev3_old = 0; // старая роль устройства для  отслеживания изменения роли на ходу
$Dev4_old = 0; // старая роль устройства для  отслеживания изменения роли на ходу
$Dev5_old = 0; // старая роль устройства для  отслеживания изменения роли на ходу

$Heat = false; // режим нагрев
$Cool = false; // режим охлаждение

// усройство заморожено, т.е его проценты не меняются *)
// в каждый момент времени только одно устройство может быть НЕ заморожено *)
// все остальные должны быть заморожены *)
$Dev1_frozen = false;
$Dev2_frozen = false;
$Dev3_frozen = false;
$Dev4_frozen = false;
$Dev5_frozen = false;

$sec_count = 0;

while (true) {
    $Work = $seg->get('Work'); // получить значение

    if ($Work == 0) {
        $Dev1_proc = 0;
        $Dev2_proc = 0;
        $Dev3_proc = 0;
        $Dev4_proc = 0;
        $Dev5_proc = 0;

        $Dev1_frozen = true;
        $Dev2_frozen = true;
        $Dev3_frozen = true;
        $Dev4_frozen = true;
        $Dev5_frozen = true;
        $seg->set('Dev1_frozen', true); // записать значение
        $seg->set('Dev2_frozen', true); // записать значение
        $seg->set('Dev3_frozen', true); // записать значение
        $seg->set('Dev4_frozen', true); // записать значение
        $seg->set('Dev5_frozen', true); // записать значение

        $Heat = false; //
        $Cool = false; //
        $seg->set('Heat', $Heat); // записать значение
        $seg->set('Cool', $Cool); // записать значение

        $Start_End = false;


    } else { // т.е. если редим работа
        $D_T = $seg->get('D_T'); // получить значение датчика температуры
        $Y_T = $seg->get('Y_T'); // получить значение уставки температуры

        $Heat = ($D_T <= $Y_T); // если меньше уставки, значит надо греть
        $Cool = ($D_T > $Y_T); // если больше уставки, значит надо охлаждать
        $seg->set('Heat', $Heat); // записать значение
        $seg->set('Cool', $Cool); // записать значение

        if (!$Start_End) {
            // определяем начальный процент
            if ($Heat) {
                // Обработка параметров для Dev1
                if ($seg->get("Dev1") == 0) {
                    $Dev1_proc = 0;
                    $Dev1_Isost = 0;
                } elseif ($seg->get("Dev1") == 1) {
                    $Dev1_proc = $seg->get("Dev1_Max");
                    $Dev1_Isost = $seg->get("Dev1_Max");
                } elseif ($seg->get("Dev1") == 2) {
                    $Dev1_proc = $seg->get("Dev1_Min");
                    $Dev1_Isost = $seg->get("Dev1_Min");
                }

                // Обработка параметров для Dev2
                if ($seg->get("Dev2") == 0) {
                    $Dev2_proc = 0;
                    $Dev2_Isost = 0;
                } elseif ($seg->get("Dev2") == 1) {
                    $Dev2_proc = $seg->get("Dev2_Max");
                    $Dev2_Isost = $seg->get("Dev2_Max");
                } elseif ($seg->get("Dev2") == 2) {
                    $Dev2_proc = $seg->get("Dev2_Min");
                    $Dev2_Isost = $seg->get("Dev2_Min");
                }

                // Обработка параметров для Dev3
                if ($seg->get("Dev3") == 0) {
                    $Dev3_proc = 0;
                    $Dev3_Isost = 0;
                } elseif ($seg->get("Dev3") == 1) {
                    $Dev3_proc = $seg->get("Dev3_Max");
                    $Dev3_Isost = $seg->get("Dev3_Max");
                } elseif ($seg->get("Dev3") == 2) {
                    $Dev3_proc = $seg->get("Dev3_Min");
                    $Dev3_Isost = $seg->get("Dev3_Min");
                }

                // Обработка параметров для Dev4
                if ($seg->get("Dev4") == 0) {
                    $Dev4_proc = 0;
                    $Dev4_Isost = 0;
                } elseif ($seg->get("Dev4") == 1) {
                    $Dev4_proc = $seg->get("Dev4_Max");
                    $Dev4_Isost = $seg->get("Dev4_Max");
                } elseif ($seg->get("Dev4") == 2) {
                    $Dev4_proc = $seg->get("Dev4_Min");
                    $Dev4_Isost = $seg->get("Dev4_Min");
                }

                // Обработка параметров для Dev5
                if ($seg->get("Dev5") == 0) {
                    $Dev5_proc = 0;
                    $Dev5_Isost = 0;
                } elseif ($seg->get("Dev5") == 1) {
                    $Dev5_proc = $seg->get("Dev5_Max");
                    $Dev5_Isost = $seg->get("Dev5_Max");
                } elseif ($seg->get("Dev5") == 2) {
                    $Dev5_proc = $seg->get("Dev5_Min");
                    $Dev5_Isost = $seg->get("Dev5_Min");
                }
            }

            if ($Cool) {
                // Обработка параметров для Dev1
                if ($seg->get("Dev1") == 0) {
                    $Dev1_proc = 0;
                    $Dev1_Isost = 0;
                } elseif ($seg->get("Dev1") == 1) {
                    $Dev1_proc = $seg->get("Dev1_Min");
                    $Dev1_Isost = $seg->get("Dev1_Min");
                } elseif ($seg->get("Dev1") == 2) {
                    $Dev1_proc = $seg->get("Dev1_Max");
                    $Dev1_Isost = $seg->get("Dev1_Max");
                }

                // Обработка параметров для Dev2
                if ($seg->get("Dev2") == 0) {
                    $Dev2_proc = 0;
                    $Dev2_Isost = 0;
                } elseif ($seg->get("Dev2") == 1) {
                    $Dev2_proc = $seg->get("Dev2_Min");
                    $Dev2_Isost = $seg->get("Dev2_Min");
                } elseif ($seg->get("Dev2") == 2) {
                    $Dev2_proc = $seg->get("Dev2_Max");
                    $Dev2_Isost = $seg->get("Dev2_Max");
                }

                // Обработка параметров для Dev3
                if ($seg->get("Dev3") == 0) {
                    $Dev3_proc = 0;
                    $Dev3_Isost = 0;
                } elseif ($seg->get("Dev3") == 1) {
                    $Dev3_proc = $seg->get("Dev3_Min");
                    $Dev3_Isost = $seg->get("Dev3_Min");
                } elseif ($seg->get("Dev3") == 2) {
                    $Dev3_proc = $seg->get("Dev3_Max");
                    $Dev3_Isost = $seg->get("Dev3_Max");
                }

                // Обработка параметров для Dev4
                if ($seg->get("Dev4") == 0) {
                    $Dev4_proc = 0;
                    $Dev4_Isost = 0;
                } elseif ($seg->get("Dev4") == 1) {
                    $Dev4_proc = $seg->get("Dev4_Min");
                    $Dev4_Isost = $seg->get("Dev4_Min");
                } elseif ($seg->get("Dev4") == 2) {
                    $Dev4_proc = $seg->get("Dev4_Max");
                    $Dev4_Isost = $seg->get("Dev4_Max");
                }

                // Обработка параметров для Dev5
                if ($seg->get("Dev5") == 0) {
                    $Dev5_proc = 0;
                    $Dev5_Isost = 0;
                } elseif ($seg->get("Dev5") == 1) {
                    $Dev5_proc = $seg->get("Dev5_Min");
                    $Dev5_Isost = $seg->get("Dev5_Min");
                } elseif ($seg->get("Dev5") == 2) {
                    $Dev5_proc = $seg->get("Dev5_Max");
                    $Dev5_Isost = $seg->get("Dev5_Max");
                }
            }


            // Выставляем начальную заморозку каждому устройству, определяем кто первый начнет действовать
            if ($Cool) {
                // Первый блок условий
                if ($seg->get("Dev1") > 0) {
                    $Dev1_frozen = false;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }

                // Второй блок условий
                if ($seg->get("Dev1") == 0 && $seg->get("Dev2") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = false;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }

                // Продолжаем аналогично для оставшихся блоков условий

                // Третий блок условий
                if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = false;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }
                // Четвертый блок условий
                if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = false;
                    $Dev5_frozen = true;
                }

                // Пятый блок условий
                if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev5") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = false;
                }

                // Шестой блок условий
                if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev5") == 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }
            }

            if ($Heat) {
                // Проверка Dev5 > 0
                if ($seg->get("Dev5") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = false;
                }

                // Проверка Dev5 == 0 и Dev4 > 0
                if ($seg->get("Dev5") == 0 && $seg->get("Dev4") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = false;
                    $Dev5_frozen = true;
                }

                // Проверка Dev5 == 0, Dev4 == 0 и Dev3 > 0
                if ($seg->get("Dev5") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev3") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = false;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }

                // Проверка Dev5 == 0, Dev4 == 0, Dev3 == 0 и Dev2 > 0
                if ($seg->get("Dev5") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev2") > 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = false;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }

                // Проверка Dev5 == 0, Dev4 == 0, Dev3 == 0 и Dev2 == 0  Dev1 > 0
                if ($seg->get("Dev5") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev1") > 0) {
                    $Dev1_frozen = false;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }

                // Проверка Dev5 == 0 и все остальные Dev также == 0
                if ($seg->get("Dev5") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev1") == 0) {
                    $Dev1_frozen = true;
                    $Dev2_frozen = true;
                    $Dev3_frozen = true;
                    $Dev4_frozen = true;
                    $Dev5_frozen = true;
                }
            }

            $Start_End = true;
        }


        // Если какое-то устройство на ходу подключилось, то его процент переинициализируется
        if ($Dev1_old != $seg->get("Dev1") && $Dev1_old == 0) {
            if ($Heat) {
                if ($seg->get("Dev1") == 1) {
                    $Dev1_proc = $seg->get("Dev1_Max");
                    ${"Dev1_Isost"} = $seg->get("Dev1_Max");
                }
                if ($seg->get("Dev1") == 2) {
                    $Dev1_proc = $seg->get("Dev1_Min");
                    ${"Dev1_Isost"} = $seg->get("Dev1_Min");
                }
            }
            if ($Cool) {
                if ($seg->get("Dev1") == 1) {
                    $Dev1_proc = $seg->get("Dev1_Min");
                    ${"Dev1_Isost"} = $seg->get("Dev1_Min");
                }
                if ($seg->get("Dev1") == 2) {
                    $Dev1_proc = $seg->get("Dev1_Max");
                    ${"Dev1_Isost"} = $seg->get("Dev1_Max");
                }
            }
            // если все устройства  кроме подключаемого не задействованы, то размораживаем подключаемое
            if ($seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev5") == 0) {
                $Dev1_frozen = false;
            }
        }
        // Проверяем переинициализацию для Dev2
        if ($Dev2_old != $seg->get("Dev2") && $Dev2_old == 0) {
            if ($Heat) {
                if ($seg->get("Dev2") == 1) {
                    $$Dev2_proc = $seg->get("Dev2_Max");
                    $Dev2_Isost = $seg->get("Dev2_Max");
                }
                if ($seg->get("Dev2") == 2) {
                    $Dev2_proc = $seg->get("Dev2_Min");
                    $Dev2_Isost = $seg->get("Dev2_Min");
                }
            }
            if ($Cool) {
                if ($seg->get("Dev2") == 1) {
                    $Dev2_proc = $seg->get("Dev2_Min");
                    $Dev2_Isost = $seg->get("Dev2_Min");
                }
                if ($seg->get("Dev2") == 2) {
                    $Dev2_proc = $seg->get("Dev2_Max");
                    $Dev2_Isost = $seg->get("Dev2_Max");
                }
            }
            // если все устройства  кроме подключаемого не задействованы, то размораживаем подключаемое
            if ($seg->get("Dev1") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev5") == 0) {
                $Dev2_frozen = false;
            }
        }

        // Проверяем переинициализацию для Dev3
        if ($Dev3_old != $seg->get("Dev3") && $Dev3_old == 0) {
            if ($Heat) {
                if ($seg->get("Dev3") == 1) {
                    $Dev3_proc = $seg->get("Dev3_Max");
                    $Dev3_Isost = $seg->get("Dev3_Max");
                }
                if ($seg->get("Dev3") == 2) {
                    $Dev3_proc = $seg->get("Dev3_Min");
                    $Dev3_Isost = $seg->get("Dev3_Min");
                }
            }
            if ($Cool) {
                if ($seg->get("Dev3") == 1) {
                    $Dev3_proc = $seg->get("Dev3_Min");
                    $Dev3_Isost = $seg->get("Dev3_Min");
                }
                if ($seg->get("Dev3") == 2) {
                    $Dev3_proc = $seg->get("Dev3_Max");
                    $Dev3_Isost = $seg->get("Dev3_Max");
                }
            }
            // если все устройства  кроме подключаемого не задействованы, то размораживаем подключаемое
            if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev4") == 0 && $seg->get("Dev5") == 0) {
                $Dev3_frozen = false;
            }
        }

        // Проверяем переинициализацию для Dev4
        if ($Dev4_old != $seg->get("Dev4") && $Dev4_old == 0) {
            if ($Heat) {
                if ($seg->get("Dev4") == 1) {
                    $Dev4_proc = $seg->get("Dev4_Max");
                    $Dev4_Isost = $seg->get("Dev4_Max");
                }
                if ($seg->get("Dev4") == 2) {
                    $Dev4_proc = $seg->get("Dev4_Min");
                    $Dev4_Isost = $seg->get("Dev4_Min");
                }
            }
            if ($Cool) {
                if ($seg->get("Dev4") == 1) {
                    $$Dev4_proc = $seg->get("Dev4_Min");
                    $Dev4_Isost = $seg->get("Dev4_Min");
                }
                if ($seg->get("Dev4") == 2) {
                    $Dev4_proc = $seg->get("Dev4_Max");
                    $Dev4_Isost = $seg->get("Dev4_Max");
                }
            }
            if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev5") == 0) {
                $Dev4_frozen = false;
            }
        }

        // Проверяем переинициализацию для Dev5
        if ($Dev5_old != $seg->get("Dev5") && $Dev5_old == 0) {
            if ($Heat) {
                if ($seg->get("Dev5") == 1) {
                    $Dev5_proc = $seg->get("Dev5_Max");
                    $Dev5_Isost = $seg->get("Dev5_Max");
                }
                if ($seg->get("Dev5") == 2) {
                    $Dev5_proc = $seg->get("Dev5_Min");
                    $Dev5_Isost = $seg->get("Dev5_Min");
                }
            }
            if ($Cool) {
                if ($seg->get("Dev5") == 1) {
                    $Dev5_proc = $seg->get("Dev5_Min");
                    $Dev5_Isost = $seg->get("Dev5_Min");
                }
                if ($seg->get("Dev5") == 2) {
                    $Dev5_proc = $seg->get("Dev5_Max");
                    $Dev5_Isost = $seg->get("Dev5_Max");
                }
            }
            if ($seg->get("Dev1") == 0 && $seg->get("Dev2") == 0 && $seg->get("Dev3") == 0 && $seg->get("Dev4") == 0) {
                $Dev5_frozen = false;
            }
        }

        // Обновляем старые значения после выполнения всех проверок и присвоений
        $Dev1_old = $seg->get("Dev1");
        $Dev2_old = $seg->get("Dev2");
        $Dev3_old = $seg->get("Dev3");
        $Dev4_old = $seg->get("Dev4");
        $Dev5_old = $seg->get("Dev5");


        // Если Dev1 устройство на ходу стало незадействованным при этом именно оно было активно (не заморожено),
        // то алгоритм должен передать управление следующей ступени или предыдущей
        if ($Dev1_frozen == false && $seg->get("Dev1") == 0) {
            $Dev1_frozen = true;
            $seg->set("Dev1_frozen", $Dev1_frozen);
            $Dev1_proc = 0;
            if ($Heat && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Heat && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Heat && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Heat && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            }

            if ($Cool && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Cool && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Cool && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Cool && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            }
        }

        // Если устройство Dev2 на ходу стало незадействованным при этом именно оно было активно (не заморожено),
        // то алгоритм должен передать управление следующей ступени или предыдущей
        if ($Dev2_frozen == false && $seg->get("Dev2") == 0) {
            $Dev2_frozen = true;
            $seg->set("Dev2_frozen", $Dev2_frozen);
            $Dev2_proc = 0;
            if ($Heat && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Heat && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Heat && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Heat && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            }

            if ($Cool && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Cool && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Cool && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Cool && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            }
        }

        // Если устройство Dev3 на ходу стало незадействованным при этом именно оно было активно (не заморожено),
        // то алгоритм должен передать управление следующей ступени или предыдущей
        if ($Dev3_frozen == false && $seg->get("Dev3") == 0) {
            $Dev3_frozen = true;
            $seg->set("Dev3_frozen", $Dev3_frozen);
            $$Dev3_proc = 0;
            if ($Heat && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Heat && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Heat && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Heat && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            }

            if ($Cool && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Cool && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Cool && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Cool && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            }
        }

        // Если устройство Dev4 на ходу стало незадействованным при этом именно оно было активно (не заморожено),
        // то алгоритм должен передать управление следующей ступени или предыдущей
        if ($Dev4_frozen == false && $seg->get("Dev4") == 0) {
            $Dev4_frozen = true;
            $seg->set("Dev4_frozen", $Dev4_frozen);
            $Dev4_proc = 0;
            if ($Heat && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            } elseif ($Heat && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Heat && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Heat && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            }

            if ($Cool && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Cool && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Cool && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Cool && $seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
                $seg->set("Dev5_frozen", $Dev5_frozen);
            }
        }

        // Если устройство Dev5 на ходу стало незадействованным при этом именно оно было активно (не заморожено),
        // то алгоритм должен передать управление следующей ступени или предыдущей
        if ($Dev5_frozen == false && $seg->get("Dev5") == 0) {
            $Dev5_frozen = true;
            $seg->set("Dev5_frozen", $Dev5_frozen);
            $Dev5_proc = 0;
            if ($Heat && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            } elseif ($Heat && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Heat && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Heat && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            }

            if ($Cool && $seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
                $seg->set("Dev4_frozen", $Dev4_frozen);
            } elseif ($Cool && $seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
                $seg->set("Dev3_frozen", $Dev3_frozen);
            } elseif ($Cool && $seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
                $seg->set("Dev2_frozen", $Dev2_frozen);
            } elseif ($Cool && $seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
                $seg->set("Dev1_frozen", $Dev1_frozen);
            }
        }


        // надо греть ПЕРВАЯ ступень разгоняется как нагреватель
        if ($Work && $Heat && !$Dev1_frozen && $Dev1_proc < $seg->get("Dev1_Max") && $seg->get("Dev1") == 1) {
            $Dev1_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev1_P");
            $Dev1_Isost = $Dev1_Isost + $Dev1_Psost / $seg->get("Dev1_I");
            $Dev1_proc = $Dev1_Psost + $Dev1_Isost;
        }
        if ($Work && $Heat && !$Dev1_frozen && $Dev1_proc >= $seg->get("Dev1_Max") && $seg->get("Dev1") == 1) {
            // Если первая ступень как нагреватель вышла на максимум, то морозим её
            $Dev1_proc = $seg->get("Dev1_Max");
            $Dev1_frozen = true;
            if ($seg->get("Dev2") > 0) { // размораживаем вторую
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev3") > 0) { // или размораживаем третью
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") > 0) { // или размораживаем четвертую
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) { // или размораживаем пятую
                $Dev5_frozen = false;
            } else {
                $Dev1_frozen = false; // или размораживаем снова первую если других нет
            }
        }

        // Надо охлаждать ПЕРВАЯ ступень снижает процент как нагреватель
        if ($Work && $Cool && !$Dev1_frozen && $Dev1_proc > $seg->get("Dev1_Min") && $seg->get("Dev1") == 1) {
            $Dev1_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev1_P");
            $Dev1_Isost = $Dev1_Isost + $Dev1_Psost / $seg->get("Dev1_I");
            $Dev1_proc = $Dev1_Psost + $Dev1_Isost;
        }
        if ($Work && $Cool && !$Dev1_frozen && $Dev1_proc <= $seg->get("Dev1_Min") && $seg->get("Dev1") == 1) {
            // Если первая ступень как нагреватель вернулась на минимум, то просто приравниваем её к минимуму и не морозим так как она первая
            $Dev1_proc = $seg->get("Dev1_Min");

            // но далее проверяем нет ли ещё ступеней из которых можно выжать холод или убрать тепло
            // такая ситуация возможна если была смена ролей
            if ($seg->get("Dev5") == 1 && $seg->get("Dev5_proc") > $seg->get("Dev5_Min")) {
                $Dev1_frozen = true;
                $Dev5_frozen = false;
            } elseif ($seg->get("Dev5") == 2 && $seg->get("Dev5_proc") < $seg->get("Dev5_Max")) {
                $Dev1_frozen = true;
                $Dev5_frozen = false;
            } elseif ($seg->get("Dev4") == 1 && $seg->get("Dev4_proc") > $seg->get("Dev4_Min")) {
                $Dev1_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev4") == 2 && $seg->get("Dev4_proc") < $seg->get("Dev4_Max")) {
                $Dev1_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev3") == 1 && $seg->get("Dev3_proc") > $seg->get("Dev3_Min")) {
                $Dev1_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev3") == 2 && $seg->get("Dev3_proc") < $seg->get("Dev3_Max")) {
                $Dev1_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") == 1 && $seg->get("Dev2_proc") > $seg->get("Dev2_Min")) {
                $Dev1_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev2") == 2 && $seg->get("Dev2_proc") < $seg->get("Dev2_Max")) {
                $Dev1_frozen = true;
                $Dev2_frozen = false;
            }
        }

        // Надо охлаждать ПЕРВАЯ ступень как охладитель увеличивает процент
        if ($Work && $Cool && !$Dev1_frozen && $Dev1_proc < $seg->get("Dev1_Max") && $seg->get("Dev1") == 2) {
            $Dev1_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev1_P") * (-1);
            $Dev1_Isost = $Dev1_Isost + $Dev1_Psost / $seg->get("Dev1_I");
            $Dev1_proc = $Dev1_Psost + $Dev1_Isost;
        }
        if ($Work && $Cool && !$Dev1_frozen && $Dev1_proc >= $seg->get("Dev1_Max") && $seg->get("Dev1") == 2) {
            // Если первая ступень как охладитель вышла на максимум, то НЕ морозим её
            $Dev1_proc = $seg->get("Dev1_Max");

            // но далее проверяем нет ли ещё ступеней из которых можно выжать холод или убрать тепло
            // такая ситуация возможна если была смена ролей
            if ($seg->get("Dev5") == 1 && $Dev1_proc > $seg->get("Dev1_Min")) {
                $Dev1_frozen = true;
                $Dev5_frozen = false;
            } elseif ($seg->get("Dev5") == 2 && $Dev1_proc < $seg->get("Dev1_Max")) {
                $Dev1_frozen = true;
                $Dev5_frozen = false;
            } elseif ($seg->get("Dev4") == 1 && $seg->get("Dev4_proc") > $seg->get("Dev1_Min")) {
                $Dev1_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev4") == 2 && $seg->get("Dev4_proc") < $seg->get("Dev4_Max")) {
                $Dev1_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev3") == 1 && $seg->get("Dev3_proc") > $seg->get("Dev3_Min")) {
                $Dev1_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev3") == 2 && $seg->get("Dev3_proc") < $seg->get("Dev3_Max")) {
                $Dev1_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") == 1 && $seg->get("Dev2_proc") > $seg->get("Dev2_Min")) {
                $Dev1_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev2") == 2 && $seg->get("Dev2_proc") < $seg->get("Dev2_Max")) {
                $Dev1_frozen = true;
                $Dev2_frozen = false;
            }
        }

        // Надо греть, ПЕРВАЯ ступень-охладитель, снижает свой процент
        if ($Work && $Heat && !$Dev1_frozen && $Dev1_proc > $seg->get("Dev1_Min") && $seg->get("Dev1") == 2) {
            $Dev1_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev1_P") * (-1);
            $Dev1_Isost = $Dev1_Isost + $Dev1_Psost / $seg->get("Dev1_I");
            $Dev1_proc = $Dev1_Psost + $Dev1_Isost;
        }
        if ($Work && $Heat && !$Dev1_frozen && $Dev1_proc <= $seg->get("Dev1_Min") && $seg->get("Dev1") == 2) {
            // Если первая ступень как охладитель вернулась на минимум
            $Dev1_proc = $seg->get("Dev1_Min");
            $Dev1_frozen = true;
            if ($seg->get("Dev2") > 0) { // размораживаем вторую
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev3") > 0) { // или размораживаем третью
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") > 0) { // или размораживаем четвертую
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) { // или размораживаем пятую
                $Dev5_frozen = false;
            } else {
                $Dev1_frozen = false; // или размораживаем снова первую если других нет
            }
        }


        // надо греть ВТОРАЯ ступень разгоняется как нагреватель
        if ($Work && $Heat && !$Dev2_frozen && $Dev2_proc < $seg->get("Dev2_Max") && $seg->get("Dev2") == 1) {
            $Dev2_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev2_P");
            $Dev2_Isost = $Dev2_Isost + $Dev2_Psost / $seg->get("Dev2_I");
            $Dev2_proc = $Dev2_Psost + $Dev2_Isost;
        }

        if ($Work && $Heat && !$Dev2_frozen && $Dev2_proc >= $seg->get("Dev2_Max") && $seg->get("Dev2") == 1) { // Если ВТОРАЯ ступень как нагреватель вышла на максимум, то морозим её
            $Dev2_proc = $seg->get("Dev2_Max");
            $Dev2_frozen = true;
            if ($seg->get("Dev3") > 0) { // или размораживаем третью
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") > 0) { // или размораживаем четвертую
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) { // или размораживаем пятую
                $Dev5_frozen = false;
            } else {
                $Dev2_frozen = false; // или размораживаем снова ВТОРАЯ если других за ней нет
            }
        }

        // надо охлаждать ВТОРАЯ ступень снижает процент как нагреватель
        if ($Work && $Cool && !$Dev2_frozen && $Dev2_proc > $seg->get("Dev2_Min") && $seg->get("Dev2") == 1) {
            $Dev2_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev2_P");
            $Dev2_Isost = $Dev2_Isost + $Dev2_Psost / $seg->get("Dev2_I");
            $Dev2_proc = $Dev2_Psost + $Dev2_Isost;
        }

        if ($Work && $Cool && !$Dev2_frozen && $Dev2_proc <= $seg->get("Dev2_Min") && $seg->get("Dev2") == 1) { // Если ВТОРАЯ ступень как нагреватель вернулась на минимум
            $Dev2_proc = $seg->get("Dev2_Min");
            $Dev2_frozen = true;
            if ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev2_frozen = false;
            }
        }

        // Надо охлаждать ВТОРАЯ ступень как охладитель увеличивает процент
        if ($Work && $Cool && !$Dev2_frozen && $Dev2_proc < $seg->get("Dev2_Max") && $seg->get("Dev2") == 2) {
            $Dev2_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev2_P") * (-1);
            $Dev2_Isost = $Dev2_Isost + $Dev2_Psost / $seg->get("Dev2_I");
            $Dev2_proc = $Dev2_Psost + $Dev2_Isost;
        }

        // Надо охлаждать ВТОРАЯ ступень как охладитель вышла на максимум
        if ($Work && $Cool && !$Dev2_frozen && $Dev2_proc >= $seg->get("Dev2_Max") && $seg->get("Dev2") == 2) {
            $Dev2_proc = $seg->get("Dev2_Max");
            $Dev2_frozen = true;
            if ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } elseif ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
            } else {
                $Dev2_frozen = false;
            }
        }

        // Надо греть, ВТОРАЯ ступень-охладитель, снижает свой процент
        if ($Work && $Heat && !$Dev2_frozen && $Dev2_proc > $seg->get("Dev2_Min") && $seg->get("Dev2") == 2) {
            $Dev2_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev2_P") * (-1);
            $Dev2_Isost = $Dev2_Isost + $Dev2_Psost / $seg->get("Dev2_I");
            $Dev2_proc = $Dev2_Psost + $Dev2_Isost;
        }
        if ($Work && $Heat && !$Dev2_frozen && $Dev2_proc <= $seg->get("Dev2_Min") && $seg->get("Dev2") == 2) {
            // Если ВТОРАЯ ступень как охладитель вышла на минимум
            $Dev2_proc = $seg->get("Dev2_Min");
            $Dev2_frozen = true;
            if ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
            } else {
                $Dev2_frozen = false; // или размораживаем снова ВТОРАЯ если других нет
            }
        }


        // Надо греть ТРЕТЬЯ ступень разгоняется как нагреватель
        if ($Work && $Heat && !$Dev3_frozen && $Dev3_proc < $seg->get("Dev3_Max") && $seg->get("Dev3") == 1) {
            $Dev3_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev3_P");
            $Dev3_Isost = $Dev3_Isost + $Dev3_Psost / $seg->get("Dev3_I");
            $Dev3_proc = $Dev3_Psost + $Dev3_Isost;
        }
        if ($Work && $Heat && !$Dev3_frozen && $Dev3_proc >= $seg->get("Dev3_Max") && $seg->get("Dev3") == 1) {
            // Если ТРЕТЬЯ ступень как нагреватель вышла на максимум, то морозим её
            $Dev3_proc = $seg->get("Dev3_Max");
            $Dev3_frozen = true;
            if ($seg->get("Dev4") > 0) { // или размораживаем четвертую
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) { // или размораживаем пятую
                $Dev5_frozen = false;
            } else {
                $Dev3_frozen = false; // или размораживаем снова ТРЕТЬЯ если других за ней нет
            }
        }

        // надо охлаждать ТРЕТЬЯ ступень снижает процент как нагреватель
        if ($Work && $Cool && !$Dev3_frozen && $Dev3_proc > $seg->get("Dev3_Min") && $seg->get("Dev3") == 1) {
            $Dev3_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev3_P");
            $Dev3_Isost = $Dev3_Isost + $Dev3_Psost / $seg->get("Dev3_I");
            $Dev3_proc = $Dev3_Psost + $Dev3_Isost;
        }
        if ($Work && $Cool && !$Dev3_frozen && $Dev3_proc <= $seg->get("Dev3_Min") && $seg->get("Dev3") == 1) {
            // Если ТРЕТЬЯ ступень как нагреватель вернулась на минимум
            $Dev3_proc = $seg->get("Dev3_Min");
            $Dev3_frozen = true;
            if ($seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev3_frozen = false;
            }
        }

        // Надо охлаждать ТРЕТЬЯ ступень как охладитель увеличивает процент
        if ($Work && $Cool && !$Dev3_frozen && $Dev3_proc < $seg->get("Dev3_Max") && $seg->get("Dev3") == 2) {
            $Dev3_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev3_P") * (-1);
            $Dev3_Isost = $Dev3_Isost + $Dev3_Psost / $seg->get("Dev3_I");
            $Dev3_proc = $Dev3_Psost + $Dev3_Isost;
        }
        if ($Work && $Cool && !$Dev3_frozen && $Dev3_proc >= $seg->get("Dev3_Max") && $seg->get("Dev3") == 2) {
            // Если ТРЕТЬЯ ступень как охладитель вышла на максимум
            $Dev3_proc = $seg->get("Dev3_Max");
            $Dev3_frozen = true;
            if ($seg->get("Dev2") > 0) { // размораживаем первую
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev3_frozen = false;
            }
        }

        // Надо греть, ТРЕТЬЯ ступень-охладитель, снижает свой процент
        if ($Work && $Heat && !$Dev3_frozen && $Dev3_proc > $seg->get("Dev3_Min") && $seg->get("Dev3") == 2) {
            $Dev3_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev3_P") * (-1);
            $Dev3_Isost = $Dev3_Isost + $Dev3_Psost / $seg->get("Dev3_I");
            $Dev3_proc = $Dev3_Psost + $Dev3_Isost;
        }
        if ($Work && $Heat && !$Dev3_frozen && $Dev3_proc <= $seg->get("Dev3_Min") && $seg->get("Dev3") == 2) {
            // Если ТРЕТЬЯ ступень как охладитель вышла на минимум
            $Dev3_proc = $seg->get("Dev3_Min");
            $Dev3_frozen = true;
            if ($seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
            } else {
                $Dev3_frozen = false; // или размораживаем снова ТРЕТЬЮ если других нет
            }
        }

        // надо греть ЧЕТВЁРТАЯ ступень разгоняется как нагреватель
        if ($Work && $Heat && !$Dev4_frozen && $Dev4_proc < $seg->get("Dev4_Max") && $seg->get("Dev4") == 1) {
            $Dev4_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev4_P");
            $Dev4_Isost = $Dev4_Isost + $Dev4_Psost / $seg->get("Dev4_I");
            $Dev4_proc = $Dev4_Psost + $Dev4_Isost;
        }
        if ($Work && $Heat && !$Dev4_frozen && $Dev4_proc >= $seg->get("Dev4_Max") && $seg->get("Dev4") == 1) {
            $Dev4_proc = $seg->get("Dev4_Max");
            $Dev4_frozen = true;
            if ($seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
            } else {
                $Dev4_frozen = false;
            }
        }

        // надо охлаждать ЧЕТВЁРТАЯ ступень снижает процент как нагреватель
        if ($Work && $Cool && !$Dev4_frozen && $Dev4_proc > $seg->get("Dev4_Min") && $seg->get("Dev4") == 1) {
            $Dev4_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev4_P");
            $Dev4_Isost = $Dev4_Isost + $Dev4_Psost / $seg->get("Dev4_I");
            $Dev4_proc = $Dev4_Psost + $Dev4_Isost;
        }
        if ($Work && $Cool && !$Dev4_frozen && $Dev4_proc <= $seg->get("Dev4_Min") && $seg->get("Dev4") == 1) {
            $Dev4_proc = $seg->get("Dev4_Min");
            $Dev4_frozen = true;
            if ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev4_frozen = false;
            }
        }

        // Надо охлаждать ЧЕТВЁРТАЯ ступень как охладитель увеличивает процент
        if ($Work && $Cool && !$Dev4_frozen && $Dev4_proc < $seg->get("Dev4_Max") && $seg->get("Dev4") == 2) {
            $Dev4_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev4_P") * (-1);
            $Dev4_Isost = $Dev4_Isost + $Dev4_Psost / $seg->get("Dev4_I");
            $Dev4_proc = $Dev4_Psost + $Dev4_Isost;
        }
        if ($Work && $Cool && !$Dev4_frozen && $Dev4_proc >= $seg->get("Dev4_Max") && $seg->get("Dev4") == 2) {
            $Dev4_proc = $seg->get("Dev4_Max");
            $Dev4_frozen = true;
            if ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev4_frozen = false;
            }
        }

        // Надо греть, ЧЕТВЁРТАЯ ступень-охладитель, снижает свой процент
        if ($Work && $Heat && !$Dev4_frozen && $Dev4_proc > $seg->get("Dev4_Min") && $seg->get("Dev4") == 2) {
            $Dev4_Psost = ($seg->get("Y_T") - $seg->get("D_T")) * $seg->get("Dev4_P") * (-1);
            $Dev4_Isost = $Dev4_Isost + $Dev4_Psost / $seg->get("Dev4_I");
            $Dev4_proc = $Dev4_Psost + $Dev4_Isost;
        }
        if ($Work && $Heat && !$Dev4_frozen && $Dev4_proc <= $seg->get("Dev4_Min") && $seg->get("Dev4") == 2) {
            $Dev4_proc = $seg->get("Dev4_Min");
            $Dev4_frozen = true;
            if ($seg->get("Dev5") > 0) {
                $Dev5_frozen = false;
            } else {
                $Dev4_frozen = false;
            }
        }


        // надо греть ПЯТАЯ ступень разгоняется как нагреватель
        if ($Work && $Heat && !$Dev5_frozen && $Dev5_proc < $seg->get("Dev5_Max") && $seg->get("Dev5") == 1) {
            $Dev5_Psost = ($Y_T - $D_T) * $seg->get("Dev5_P");
            $Dev5_Isost += $Dev5_Psost / $seg->get("Dev5_I");
            $Dev5_proc = $Dev5_Psost + $Dev5_Isost;
        }

        // надо охлаждать ПЯТАЯ ступень снижает процент как нагреватель
        if ($Work && $Cool && !$Dev5_frozen && $Dev5_proc > $seg->get("Dev5_Min") && $seg->get("Dev5") == 1) {
            $Dev5_Psost = ($Y_T - $D_T) * $seg->get("Dev5_P");
            $Dev5_Isost += $Dev5_Psost / $seg->get("Dev5_I");
            $Dev5_proc = $Dev5_Psost + $Dev5_Isost;
        }

        // не морозим её, так как она последняя, но далее проверяем, нет ли ещё ступеней, из которых можно выжать тепло или убрать холод
        if ($Work && $Heat && !$Dev5_frozen && $Dev5_proc >= $seg->get("Dev5_Max") && $seg->get("Dev5") == 1) {
            $Dev5_proc = $seg->get("Dev5_Max");

            if ($seg->get("Dev1") == 1 && $Dev1_proc < $seg->get("Dev1_Max")) {
                $Dev5_frozen = true;
                $Dev1_frozen = false;
            } elseif ($seg->get("Dev1") == 2 && $Dev1_proc > $seg->get("Dev1_Min")) {
                $Dev5_frozen = true;
                $Dev1_frozen = false;
            } elseif ($seg->get("Dev2") == 1 && $Dev2_proc < $seg->get("Dev2_Max")) {
                $Dev5_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev2") == 2 && $Dev2_proc > $seg->get("Dev2_Min")) {
                $Dev5_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev3") == 1 && $Dev3_proc < $seg->get("Dev3_Max")) {
                $Dev5_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev3") == 2 && $Dev3_proc > $seg->get("Dev3_Min")) {
                $Dev5_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") == 1 && $Dev4_proc < $seg->get("Dev4_Max")) {
                $Dev5_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev4") == 2 && $Dev4_proc > $seg->get("Dev4_Min")) {
                $Dev5_frozen = true;
                $Dev4_frozen = false;
            }
        }

        // не морозим её , так как она последняя, но далее проверяем, нет ли ещё ступеней, из которых можно выжать тепло или убрать холод
        if ($Work && $Cool && !$Dev5_frozen && $Dev5_proc <= $seg->get("Dev5_Min") && $seg->get("Dev5") == 1) {
            $Dev5_proc = $seg->get("Dev5_Min");
            $Dev5_frozen = true;

            if ($seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev5_frozen = false;
            }
        }

        // Надо охлаждать ПЯТАЯ ступень как охладитель увеличивает процент
        if ($Work && $Cool && !$Dev5_frozen && $Dev5_proc < $seg->get("Dev5_Max") && $seg->get("Dev5") == 2) {
            $Dev5_Psost = ($Y_T - $D_T) * $seg->get("Dev5_P") * (-1);
            $Dev5_Isost += $Dev5_Psost / $seg->get("Dev5_I");
            $Dev5_proc = $Dev5_Psost + $Dev5_Isost;
        }

        if ($Work && $Cool && !$Dev5_frozen && $Dev5_proc >= $seg->get("Dev5_Max") && $seg->get("Dev5") == 2) {
            $Dev5_proc = $seg->get("Dev5_Max");
            $Dev5_frozen = true;

            if ($seg->get("Dev4") > 0) {
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev3") > 0) {
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev2") > 0) {
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev1") > 0) {
                $Dev1_frozen = false;
            } else {
                $Dev5_frozen = false;
            }
        }

        // Надо греть, ПЯТАЯ ступень-охладитель, снижает свой процент
        if ($Work && $Heat && !$Dev5_frozen && $Dev5_proc > $seg->get("Dev5_Min") && $seg->get("Dev5") == 2) {
            $Dev5_Psost = ($Y_T - $D_T) * $seg->get("Dev5_P") * (-1);
            $Dev5_Isost += $Dev5_Psost / $seg->get("Dev5_I");
            $Dev5_proc = $Dev5_Psost + $Dev5_Isost;
        }

        if ($Work && $Heat && !$Dev5_frozen && $Dev5_proc <= $seg->get("Dev5_Min") && $seg->get("Dev5") == 2) {
            $Dev5_proc = $seg->get("Dev5_Min");

            if ($seg->get("Dev1") == 1 && $Dev1_proc < $seg->get("Dev1_Max")) {
                $Dev5_frozen = true;
                $Dev1_frozen = false;
            } elseif ($seg->get("Dev1") == 2 && $Dev1_proc > $seg->get("Dev1_Min")) {
                $Dev5_frozen = true;
                $Dev1_frozen = false;
            } elseif ($seg->get("Dev2") == 1 && $Dev2_proc < $seg->get("Dev2_Max")) {
                $Dev5_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev2") == 2 && $Dev2_proc > $seg->get("Dev2_Min")) {
                $Dev5_frozen = true;
                $Dev2_frozen = false;
            } elseif ($seg->get("Dev3") == 1 && $Dev3_proc < $seg->get("Dev3_Max")) {
                $Dev5_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev3") == 2 && $Dev3_proc > $seg->get("Dev3_Min")) {
                $Dev5_frozen = true;
                $Dev3_frozen = false;
            } elseif ($seg->get("Dev4") == 1 && $Dev4_proc < $seg->get("Dev4_Max")) {
                $Dev5_frozen = true;
                $Dev4_frozen = false;
            } elseif ($seg->get("Dev4") == 2 && $Dev4_proc > $seg->get("Dev4_Min")) {
                $Dev5_frozen = true;
                $Dev4_frozen = false;
            }
        }

        // обнуляем незадействованные устройства
        if ($seg->get("Dev1") == 0) {
            $Dev1_proc = 0;
        }
        if ($seg->get("Dev2") == 0) {
            $Dev2_proc = 0;
        }
        if ($seg->get("Dev3") == 0) {
            $Dev3_proc = 0;
        }
        if ($seg->get("Dev4") == 0) {
            $Dev4_proc = 0;
        }
        if ($seg->get("Dev5") == 0) {
            $Dev5_proc = 0;
        }

        // подрезаем проценты снизу и сверху, если вдруг они менялись на ходу
        if ($Dev1_proc > $seg->get("Dev1_Max")) {
            $Dev1_proc = $seg->get("Dev1_Max");
        }
        if ($Dev2_proc > $seg->get("Dev2_Max")) {
            $Dev2_proc = $seg->get("Dev2_Max");
        }
        if ($Dev3_proc > $seg->get("Dev3_Max")) {
            $Dev3_proc = $seg->get("Dev3_Max");
        }
        if ($Dev4_proc > $seg->get("Dev4_Max")) {
            $Dev4_proc = $seg->get("Dev4_Max");
        }
        if ($Dev5_proc > $seg->get("Dev5_Max")) {
            $Dev5_proc = $seg->get("Dev5_Max");
        }

        if ($Dev1_proc < $seg->get("Dev1_Min")) {
            $Dev1_proc = $seg->get("Dev1_Min");
        }
        if ($Dev2_proc < $seg->get("Dev2_Min")) {
            $Dev2_proc = $seg->get("Dev2_Min");
        }
        if ($Dev3_proc < $seg->get("Dev3_Min")) {
            $Dev3_proc = $seg->get("Dev3_Min");
        }
        if ($Dev4_proc < $seg->get("Dev4_Min")) {
            $Dev4_proc = $seg->get("Dev4_Min");
        }
        if ($Dev5_proc < $seg->get("Dev5_Min")) {
            $Dev5_proc = $seg->get("Dev5_Min");
        }
    }

    $seg->set("Dev1_frozen", $Dev1_frozen);
    $seg->set("Dev2_frozen", $Dev2_frozen);
    $seg->set("Dev3_frozen", $Dev3_frozen);
    $seg->set("Dev4_frozen", $Dev4_frozen);
    $seg->set("Dev5_frozen", $Dev5_frozen);

    $seg->set('Dev1_proc', $Dev1_proc); // записать значение
    $seg->set('Dev2_proc', $Dev2_proc); // записать значение
    $seg->set('Dev3_proc', $Dev3_proc); // записать значение
    $seg->set('Dev4_proc', $Dev4_proc); // записать значение
    $seg->set('Dev5_proc', $Dev5_proc); // записать значение

    $sec_count++;

    $seg->set('count', $sec_count); // записать значение
    usleep(1000000); // пауза на 1 секунду

    // пишем счётчик в консоль чтоб просто было видно что алгоритм работает
    echo($sec_count . "\n");
    echo("\n");

}

?>