<?php

class SegShm {
    private $ctx = NULL;
    
	function __construct() {
		$this->cfg = parse_ini_file('/projects/load_files.srv', true, INI_SCANNER_RAW);
		
		$memSize = intval($this->cfg['Slave']['ShmSize']);
		
		$shmKey = ftok('/dev/shm/wsi', 'b');
		$this->shm = shmop_open($shmKey, 'w', 0644, $memSize);
		
		$varTypeList = [ 'Coil', 'Instat', 'Inreg', 'Holdreg' ];
		
		foreach ($varTypeList as $vt) {
			if (!isset($this->cfg[$vt])) continue;
			
			$varList = $this->cfg[$vt];
			foreach ($varList as $var) {
				$var = iconv('windows-1251', 'utf-8', $var);
				
				list($tmp, $vSize, $vType, $vOff, $vSt, $vFin, $vName)
					= split(",", $var);
				
				$vd['type'] = intval($vType);
				//$vd['addr'] = hexdec($vOff) + (hexdec($vFin) - hexdec($vSt));
				$vd['addr'] = hexdec($vOff) + 48; // sizeof(pthread_mutex_t)=24 for 32bit, =48 for 64bit;
				$this->vd[$vName] = $vd;
			}
		}
	}

    public function get($name) {
		$vd = $this->vd[$name];
		
		switch ($vd['type']) {
		case 0: // булевый
			$res = shmop_read($this->shm, $vd['addr'], 1);
			$d = unpack('C1', $res);
			return $d[1] == 1;
		case 1: // целый двухбайтовый
			$res = shmop_read($this->shm, $vd['addr'], 2);
			$d = unpack('s1', $res);
			return $d[1];
		case 3: // с плавающей точкой
			$res = shmop_read($this->shm, $vd['addr'], 4);
			$d = unpack('f1', $res);
			return $d[1];
		default:
			error_log('Unknown data type');
		}
	}
	
	public function set($name, $val) {
		$vd = $this->vd[$name];
		
		switch ($vd['type']) {
		case 0: // булевый
			if ($val) $v = 1;
			else $v = 0;
			
			$d = pack('C1', $v);
			break;
		case 1: // целый двухбайтовый
			$d = pack('s1', $val);
			break;
		case 3: // с плавающей точкой
			$d = pack('f1', $val);
			break;
		default:
			error_log('Unknown data type');
		}
		
		shmop_write($this->shm, $d, $vd['addr']);
	}
	
	public function dump() {
		echo('<table border="1">');
		
		foreach ($this->vd as $name => $vd) {
			echo('<tr><td>' . $name . '</td><td>' . $this->get($name) . '</td></tr>');
		}
		
		echo('</table>');
	}
}

?>







